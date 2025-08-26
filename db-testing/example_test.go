package dbtesting

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// Repository example for testing
type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) GetByID(id uint) (*User, error) {
	var user User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *UserRepository) GetByName(name string) (*User, error) {
	var user User
	err := r.db.Where("name = ?", name).First(&user).Error
	return &user, err
}

// TestRepositoryExample shows how to test repositories
func TestRepositoryExample(t *testing.T) {
	db := CreateTestDB(t)
	err := db.AutoMigrate(&User{})
	require.NoError(t, err)

	repo := NewUserRepository(db)

	t.Run("Create and retrieve user", func(t *testing.T) {
		user := &User{Name: "Alice"}
		err := repo.Create(user)
		require.NoError(t, err)
		assert.NotZero(t, user.ID)

		found, err := repo.GetByID(user.ID)
		require.NoError(t, err)
		assert.Equal(t, "Alice", found.Name)
	})

	t.Run("Find by name", func(t *testing.T) {
		user := &User{Name: "Bob"}
		err := repo.Create(user)
		require.NoError(t, err)

		found, err := repo.GetByName("Bob")
		require.NoError(t, err)
		assert.Equal(t, user.ID, found.ID)
	})
}
