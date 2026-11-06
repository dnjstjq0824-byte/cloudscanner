package identity

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlugify(t *testing.T) {
	assert.Equal(t, "acme-corp", slugify("Acme Corp"))
	assert.Equal(t, "johns-organization", slugify("John's Organization"))
	assert.Equal(t, "a-b-c", slugify("  A__B  C "))
	assert.Equal(t, "", slugify("!!!"))
	assert.LessOrEqual(t, len(slugify("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")), 70)
}

func TestMapRoleCode(t *testing.T) {
	assert.Equal(t, RoleOwner, mapRoleCode("ORG_OWNER"))
	assert.Equal(t, RoleAdmin, mapRoleCode("ORG_ADMIN"))
	assert.Equal(t, RoleAdmin, mapRoleCode("SUPER_ADMIN"))
	assert.Equal(t, RoleViewer, mapRoleCode("ORG_VIEWER"))
	assert.Equal(t, RoleMember, mapRoleCode("CLOUD_OPERATOR"))
	assert.Equal(t, RoleMember, mapRoleCode("BUDGET_APPROVER"))
	assert.Equal(t, RoleMember, mapRoleCode("UNKNOWN_CODE"))
}

func TestIsUniqueViolation(t *testing.T) {
	assert.False(t, isUniqueViolation(nil))
	assert.True(t, isUniqueViolation(errors.New(`ERROR: duplicate key value violates unique constraint "users_email_key" (SQLSTATE 23505)`)))
	assert.False(t, isUniqueViolation(errors.New("connection refused")))
}
