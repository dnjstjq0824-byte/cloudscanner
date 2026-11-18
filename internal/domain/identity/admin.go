package identity

import (
	"context"
	"errors"
	"sort"
	"strings"
	"time"
)

// Account statuses enforced by UpdateUserStatus.
const (
	StatusPending   = "PENDING_VERIFICATION"
	StatusActive    = "ACTIVE"
	StatusSuspended = "SUSPENDED"
)

// AdminUserView joins a user with their primary organization and role.
type AdminUserView struct {
	ID               string `json:"id"`
	Email            string `json:"email"`
	FullName         string `json:"full_name"`
	Status           string `json:"status"`
	OrganizationID   string `json:"organization_id"`
	OrganizationName string `json:"organization_name"`
	Role             string `json:"role"`
	CreatedAt        string `json:"created_at"`
}

// NormalizeRole validates a token role rank.
func NormalizeRole(role string) (string, error) {
	switch strings.ToUpper(strings.TrimSpace(role)) {
	case RoleOwner:
		return RoleOwner, nil
	case RoleAdmin:
		return RoleAdmin, nil
	case RoleMember:
		return RoleMember, nil
	case RoleViewer:
		return RoleViewer, nil
	default:
		return "", errors.New("unknown role: want OWNER, ADMIN, MEMBER, or VIEWER")
	}
}

// NormalizeStatus validates an account status transition target.
func NormalizeStatus(status string) (string, error) {
	switch strings.ToUpper(strings.TrimSpace(status)) {
	case StatusActive:
		return StatusActive, nil
	case StatusSuspended:
		return StatusSuspended, nil
	case StatusPending:
		return StatusPending, nil
	default:
		return "", errors.New("unknown status: want ACTIVE, SUSPENDED, or PENDING_VERIFICATION")
	}
}

// ListUsers returns newest-first user views with limit/offset pagination.
func (s *MemoryUserStore) ListUsers(ctx context.Context, limit, offset int, status string) ([]AdminUserView, int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var views []AdminUserView
	for _, u := range s.byID {
		if status != "" && !strings.EqualFold(u.Status, status) {
			continue
		}
		orgID, orgName, role := "", "", RoleMember
		for _, org := range s.orgs {
			if org.OwnerUserID == u.ID {
				orgID, orgName = org.ID, org.Name
				role = RoleOwner
				break
			}
		}
		if r, ok := s.roles[u.ID]; ok {
			role = r
		}
		views = append(views, AdminUserView{
			ID: u.ID, Email: u.Email, FullName: u.FullName, Status: u.Status,
			OrganizationID: orgID, OrganizationName: orgName, Role: role,
			CreatedAt: u.CreatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
		})
	}
	sort.Slice(views, func(i, j int) bool { return views[i].CreatedAt > views[j].CreatedAt })
	total := len(views)
	if offset < 0 {
		offset = 0
	}
	if offset >= total {
		return []AdminUserView{}, total, nil
	}
	end := offset + limit
	if limit <= 0 || end > total {
		end = total
	}
	return views[offset:end], total, nil
}

// UpdateUserStatus transitions an account status.
func (s *MemoryUserStore) UpdateUserStatus(ctx context.Context, id, status string) error {
	st, err := NormalizeStatus(status)
	if err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	u, exists := s.byID[id]
	if !exists {
		return errors.New("user not found")
	}
	u.Status = st
	u.UpdatedAt = time.Now().UTC()
	return nil
}

// UpdateUserRole reassigns the primary organization role.
func (s *MemoryUserStore) UpdateUserRole(ctx context.Context, userID, role string) error {
	r, err := NormalizeRole(role)
	if err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.byID[userID]; !exists {
		return errors.New("user not found")
	}
	if s.roles == nil {
		s.roles = make(map[string]string)
	}
	s.roles[userID] = r
	return nil
}

// CountUsers returns the total registered user count.
func (s *MemoryUserStore) CountUsers(ctx context.Context) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.byID), nil
}
