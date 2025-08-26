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

// TestRepositoryExample shows how to test repositories with different environments
func TestRepositoryExample(t *testing.T) {
	t.Run("EnvTest with transaction isolation", func(t *testing.T) {
		db := CreateTestDB(t, EnvTest)
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
	})

	t.Run("EnvTest with debug disabled for clean output", func(t *testing.T) {
		db := CreateTestDB(t, EnvTest, DBDebugOff)
		err := db.AutoMigrate(&User{})
		require.NoError(t, err)

		repo := NewUserRepository(db)

		user := &User{Name: "Charlie"}
		err = repo.Create(user)
		require.NoError(t, err)
		assert.NotZero(t, user.ID)

		found, err := repo.GetByName("Charlie")
		require.NoError(t, err)
		assert.Equal(t, user.ID, found.ID)
	})

	t.Run("EnvDev for integration testing (may skip)", func(t *testing.T) {
		db := CreateTestDB(t, EnvDev, DBDebugOff)
		if db == nil {
			t.Skip("Development database not available")
			return
		}

		err := db.AutoMigrate(&User{})
		require.NoError(t, err)

		repo := NewUserRepository(db)

		user := &User{Name: "Dev User"}
		err = repo.Create(user)
		require.NoError(t, err)
		assert.NotZero(t, user.ID)
	})

	t.Run("Backwards compatibility", func(t *testing.T) {
		// Legacy API still works
		db := CreateTestDBLegacy(t)
		err := db.AutoMigrate(&User{})
		require.NoError(t, err)

		repo := NewUserRepository(db)

		user := &User{Name: "Legacy Alice"}
		err = repo.Create(user)
		require.NoError(t, err)
		assert.NotZero(t, user.ID)
	})
}
