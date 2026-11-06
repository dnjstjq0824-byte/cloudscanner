package identity

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cloudscanr/cloudscanr/internal/platform/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// PGUserStore implements UserStore against PostgreSQL 16 (migration 000002).
type PGUserStore struct {
	db *database.DBPool
}

// NewPGUserStore constructs a PostgreSQL-backed user repository.
func NewPGUserStore(db *database.DBPool) *PGUserStore {
	return &PGUserStore{db: db}
}

// CreateUser inserts a user, their default organization, and an ORG_OWNER
// membership atomically. Slug collisions are resolved with a numeric suffix.
func (s *PGUserStore) CreateUser(ctx context.Context, u *User, orgName string) (*User, *Organization, error) {
	emailLower := strings.ToLower(strings.TrimSpace(u.Email))
	if emailLower == "" {
		return nil, nil, errors.New("email is required")
	}
	displayName := strings.TrimSpace(u.FullName)
	if displayName == "" {
		displayName = "CloudScanr User"
	}
	status := u.Status
	if status == "" {
		status = "ACTIVE"
	}

	baseSlug := slugify(orgName)
	if baseSlug == "" {
		baseSlug = "personal-organization"
	}

	tx, err := s.db.RawPool().Begin(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("begin create-user transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	userID := uuid.NewString()
	now := time.Now().UTC()
	_, err = tx.Exec(ctx, `
		INSERT INTO users (id, email, password_hash, display_name, account_status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $6)`,
		userID, emailLower, u.PasswordHash, displayName, status, now)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, nil, errors.New("email address is already registered")
		}
		return nil, nil, fmt.Errorf("insert user: %w", err)
	}

	slug := baseSlug
	for attempt := 0; attempt < 10; attempt++ {
		candidate := slug
		if attempt > 0 {
			candidate = fmt.Sprintf("%s-%d", slug, attempt+1)
		}
		orgID := uuid.NewString()
		_, err = tx.Exec(ctx, `
			INSERT INTO organizations (id, slug, name, owner_user_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $5)`,
			orgID, candidate, orgName, userID, now)
		if err == nil {
			var roleID string
			err = tx.QueryRow(ctx, `SELECT id FROM roles WHERE code = 'ORG_OWNER'`).Scan(&roleID)
			if err != nil {
				return nil, nil, fmt.Errorf("lookup ORG_OWNER role (run seed migration 000011): %w", err)
			}
			_, err = tx.Exec(ctx, `
				INSERT INTO organization_members (organization_id, user_id, role_id, can_connect_cloud_account, can_provision_server, can_approve_server_deletion, can_approve_budget)
				VALUES ($1, $2, $3, TRUE, TRUE, TRUE, TRUE)`,
				orgID, userID, roleID)
			if err != nil {
				return nil, nil, fmt.Errorf("insert owner membership: %w", err)
			}
			if err := tx.Commit(ctx); err != nil {
				return nil, nil, fmt.Errorf("commit create-user transaction: %w", err)
			}
			return &User{
					ID:           userID,
					Email:        emailLower,
					PasswordHash: u.PasswordHash,
					FullName:     displayName,
					Status:       status,
					CreatedAt:    now,
					UpdatedAt:    now,
				}, &Organization{
					ID:          orgID,
					Name:        orgName,
					Slug:        candidate,
					OwnerUserID: userID,
					CreatedAt:   now,
				}, nil
		}
		if !isUniqueViolation(err) {
			return nil, nil, fmt.Errorf("insert organization: %w", err)
		}
	}
	return nil, nil, errors.New("could not allocate a unique organization slug")
}

// FindByEmail loads a user by lowercase email.
func (s *PGUserStore) FindByEmail(ctx context.Context, email string) (*User, error) {
	return s.findUser(ctx, `SELECT id, email, password_hash, display_name, account_status, created_at, updated_at
		FROM users WHERE email = $1`, strings.ToLower(strings.TrimSpace(email)))
}

// FindByID loads a user by UUID.
func (s *PGUserStore) FindByID(ctx context.Context, id string) (*User, error) {
	return s.findUser(ctx, `SELECT id, email, password_hash, display_name, account_status, created_at, updated_at
		FROM users WHERE id = $1`, id)
}

func (s *PGUserStore) findUser(ctx context.Context, query, arg string) (*User, error) {
	var u User
	err := s.db.RawPool().QueryRow(ctx, query, arg).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.FullName, &u.Status, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &u, nil
}

// FindOrgByOwner returns the user's primary organization with its mapped role.
func (s *PGUserStore) FindOrgByOwner(ctx context.Context, userID string) (*Organization, string, error) {
	var org Organization
	var roleCode string
	err := s.db.RawPool().QueryRow(ctx, `
		SELECT o.id, o.slug, o.name, o.owner_user_id, o.created_at, r.code
		FROM organizations o
		JOIN organization_members m ON m.organization_id = o.id AND m.user_id = $1
		JOIN roles r ON r.id = m.role_id
		ORDER BY o.created_at ASC
		LIMIT 1`, userID).Scan(&org.ID, &org.Slug, &org.Name, &org.OwnerUserID, &org.CreatedAt, &roleCode)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, "", errors.New("organization not found for user")
		}
		return nil, "", err
	}
	return &org, mapRoleCode(roleCode), nil
}

// mapRoleCode translates DB role codes to token role ranks.
func mapRoleCode(code string) string {
	switch code {
	case "ORG_OWNER":
		return RoleOwner
	case "ORG_ADMIN", "SUPER_ADMIN":
		return RoleAdmin
	case "ORG_VIEWER":
		return RoleViewer
	default:
		return RoleMember
	}
}

// slugify renders an organization name as a URL-safe slug.
func slugify(name string) string {
	lower := strings.ToLower(strings.TrimSpace(name))
	var b strings.Builder
	prevDash := false
	for _, r := range lower {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
			prevDash = false
		case r == '\'':
			// Drop apostrophes: "John's" -> "johns".
		case r == ' ' || r == '-' || r == '_':
			if !prevDash && b.Len() > 0 {
				b.WriteByte('-')
				prevDash = true
			}
		}
	}
	slug := strings.Trim(b.String(), "-")
	if len(slug) > 70 {
		slug = slug[:70]
	}
	return slug
}

// isUniqueViolation reports PostgreSQL unique-constraint violations (SQLSTATE 23505).
func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "SQLSTATE 23505") ||
		strings.Contains(err.Error(), "duplicate key")
}
