package identity

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

// roleCodeForToken maps token ranks to seeded role codes (migration 000011).
func roleCodeForToken(role string) string {
	switch role {
	case RoleOwner:
		return "ORG_OWNER"
	case RoleAdmin:
		return "ORG_ADMIN"
	case RoleViewer:
		return "ORG_VIEWER"
	default:
		return "CLOUD_OPERATOR"
	}
}

// ListUsers returns newest-first user views with limit/offset pagination.
func (s *PGUserStore) ListUsers(ctx context.Context, limit, offset int, status string) ([]AdminUserView, int, error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM users`
	var err error
	if status != "" {
		err = s.db.RawPool().QueryRow(ctx, countQuery+` WHERE account_status = $1`, strings.ToUpper(status)).Scan(&total)
	} else {
		err = s.db.RawPool().QueryRow(ctx, countQuery).Scan(&total)
	}
	if err != nil {
		return nil, 0, err
	}
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	query := `
		SELECT u.id::text, u.email, u.display_name, u.account_status, u.created_at,
			COALESCE(m.organization_id::text, ''), COALESCE(m.org_name, ''), COALESCE(m.role_code, '')
		FROM users u
		LEFT JOIN LATERAL (
			SELECT o.id AS organization_id, o.name AS org_name, r.code AS role_code
			FROM organization_members om
			JOIN organizations o ON o.id = om.organization_id
			JOIN roles r ON r.id = om.role_id
			WHERE om.user_id = u.id
			ORDER BY o.created_at ASC
			LIMIT 1
		) m ON true`
	var rows pgx.Rows
	if status != "" {
		rows, err = s.db.RawPool().Query(ctx, query+` WHERE u.account_status = $1 ORDER BY u.created_at DESC LIMIT $2 OFFSET $3`,
			strings.ToUpper(status), limit, offset)
	} else {
		rows, err = s.db.RawPool().Query(ctx, query+` ORDER BY u.created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	}
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	views := []AdminUserView{}
	for rows.Next() {
		var v AdminUserView
		var created time.Time
		var roleCode string
		if err := rows.Scan(&v.ID, &v.Email, &v.FullName, &v.Status, &created,
			&v.OrganizationID, &v.OrganizationName, &roleCode); err != nil {
			return nil, 0, err
		}
		v.Role = mapRoleCode(roleCode)
		v.CreatedAt = created.UTC().Format("2006-01-02T15:04:05Z07:00")
		views = append(views, v)
	}
	return views, total, rows.Err()
}

// UpdateUserStatus transitions an account status.
func (s *PGUserStore) UpdateUserStatus(ctx context.Context, id, status string) error {
	st, err := NormalizeStatus(status)
	if err != nil {
		return err
	}
	tag, err := s.db.RawPool().Exec(ctx,
		`UPDATE users SET account_status = $2, updated_at = $3 WHERE id = $1`,
		id, st, time.Now().UTC())
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.New("user not found")
	}
	return nil
}

// UpdateUserRole reassigns the earliest-organization membership role.
func (s *PGUserStore) UpdateUserRole(ctx context.Context, userID, role string) error {
	r, err := NormalizeRole(role)
	if err != nil {
		return err
	}
	tag, err := s.db.RawPool().Exec(ctx, `
		UPDATE organization_members om SET role_id = (SELECT id FROM roles WHERE code = $2)
		WHERE om.user_id = $1 AND om.organization_id = (
			SELECT o.id FROM organizations o
			JOIN organization_members m ON m.organization_id = o.id AND m.user_id = $1
			ORDER BY o.created_at ASC LIMIT 1
		)`, userID, roleCodeForToken(r))
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.New("user membership not found")
	}
	return nil
}

// CountUsers returns the total registered user count.
func (s *PGUserStore) CountUsers(ctx context.Context) (int, error) {
	var total int
	if err := s.db.RawPool().QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}
