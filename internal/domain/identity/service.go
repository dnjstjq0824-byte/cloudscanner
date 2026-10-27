package identity

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
)

// Argon2id parameters specified in Article 15 and WO-CS-0102
const (
	Argon2Memory      = 64 * 1024 // 64MB
	Argon2Time        = 3         // 3 iterations
	Argon2Parallelism = 4         // 4 threads
	Argon2SaltLength  = 16        // 16 bytes
	Argon2KeyLength   = 32        // 32 bytes
)

// User represents a CloudScanr platform user.
type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	FullName     string    `json:"full_name"`
	Status       string    `json:"status"` // PENDING_VERIFICATION, ACTIVE, SUSPENDED
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Organization represents a multi-tenant isolation unit.
type Organization struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	OwnerUserID string    `json:"owner_user_id"`
	CreatedAt   time.Time `json:"created_at"`
}

// UserStore defines repository contract for identity.
type UserStore interface {
	CreateUser(ctx context.Context, u *User, orgName string) (*User, *Organization, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id string) (*User, error)
}

// MemoryUserStore provides an in-memory thread-safe user repository.
type MemoryUserStore struct {
	mu    sync.RWMutex
	users map[string]*User         // email -> User
	byID  map[string]*User         // id -> User
	orgs  map[string]*Organization // orgID -> Org
}

// NewMemoryUserStore constructs a memory store.
func NewMemoryUserStore() *MemoryUserStore {
	return &MemoryUserStore{
		users: make(map[string]*User),
		byID:  make(map[string]*User),
		orgs:  make(map[string]*Organization),
	}
}

func (s *MemoryUserStore) CreateUser(ctx context.Context, u *User, orgName string) (*User, *Organization, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	emailLower := strings.ToLower(u.Email)
	if _, exists := s.users[emailLower]; exists {
		return nil, nil, errors.New("email address is already registered")
	}

	if u.ID == "" {
		u.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	u.CreatedAt = now
	u.UpdatedAt = now
	u.Email = emailLower

	org := &Organization{
		ID:          uuid.NewString(),
		Name:        orgName,
		Slug:        strings.ToLower(strings.ReplaceAll(orgName, " ", "-")),
		OwnerUserID: u.ID,
		CreatedAt:   now,
	}

	s.users[emailLower] = u
	s.byID[u.ID] = u
	s.orgs[org.ID] = org

	return u, org, nil
}

func (s *MemoryUserStore) FindByEmail(ctx context.Context, email string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, exists := s.users[strings.ToLower(email)]
	if !exists {
		return nil, errors.New("user not found")
	}
	return u, nil
}

func (s *MemoryUserStore) FindByID(ctx context.Context, id string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, exists := s.byID[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return u, nil
}

// HashPassword encodes a raw password with Argon2id.
func HashPassword(password string) (string, error) {
	if len(password) < 8 {
		return "", errors.New("password must be at least 8 characters")
	}

	salt := make([]byte, Argon2SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, Argon2Time, Argon2Memory, Argon2Parallelism, Argon2KeyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, Argon2Memory, Argon2Time, Argon2Parallelism, b64Salt, b64Hash)

	return encoded, nil
}

// VerifyPassword checks plaintext password against encoded Argon2id hash.
func VerifyPassword(password, encodedHash string) (bool, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return false, errors.New("invalid or unsupported password hash format")
	}

	var version int
	_, err := fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil {
		return false, err
	}

	var memory uint32
	var iterations uint32
	var parallelism uint8
	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	computedHash := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, uint32(len(expectedHash)))

	if subtle.ConstantTimeCompare(computedHash, expectedHash) == 1 {
		return true, nil
	}
	return false, nil
}

// ValidateEmail normalizes and checks RFC 5322 compliance.
func ValidateEmail(email string) (string, error) {
	trimmed := strings.TrimSpace(email)
	addr, err := mail.ParseAddress(trimmed)
	if err != nil {
		return "", fmt.Errorf("invalid email address: %w", err)
	}
	return strings.ToLower(addr.Address), nil
}

// Service manages user identity lifecycle.
type Service struct {
	store UserStore
}

// NewService constructs identity service.
func NewService(store UserStore) *Service {
	if store == nil {
		store = NewMemoryUserStore()
	}
	return &Service{store: store}
}

// Register registers a new user and creates their default organization.
func (s *Service) Register(ctx context.Context, email, password, fullName string) (*User, *Organization, error) {
	validEmail, err := ValidateEmail(email)
	if err != nil {
		return nil, nil, err
	}

	if len(password) < 8 {
		return nil, nil, errors.New("password must be at least 8 characters")
	}

	hash, err := HashPassword(password)
	if err != nil {
		return nil, nil, fmt.Errorf("password hashing failed: %w", err)
	}

	user := &User{
		Email:        validEmail,
		PasswordHash: hash,
		FullName:     strings.TrimSpace(fullName),
		Status:       "ACTIVE",
	}

	orgName := fmt.Sprintf("%s's Organization", user.FullName)
	if user.FullName == "" {
		orgName = "Personal Organization"
	}

	return s.store.CreateUser(ctx, user, orgName)
}

// Authenticate verifies credentials and returns user on success.
func (s *Service) Authenticate(ctx context.Context, email, password string) (*User, error) {
	validEmail, err := ValidateEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	u, err := s.store.FindByEmail(ctx, validEmail)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if u.Status == "SUSPENDED" {
		return nil, errors.New("account is suspended")
	}

	valid, err := VerifyPassword(password, u.PasswordHash)
	if err != nil || !valid {
		return nil, errors.New("invalid credentials")
	}

	return u, nil
}
