package dbtesting

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type User struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"not null"`
}

func TestCreateTestDB(t *testing.T) {
	db := CreateTestDB(t)

	// Test database works
	err := db.AutoMigrate(&User{})
	require.NoError(t, err)

	user := User{Name: "Test User"}
	err = db.Create(&user).Error
	require.NoError(t, err)
	assert.NotZero(t, user.ID)

	// Test can query back
	var found User
	err = db.First(&found, user.ID).Error
	require.NoError(t, err)
	assert.Equal(t, "Test User", found.Name)
}

func TestCreateTestDBWithTx(t *testing.T) {
	t.Run("Transaction isolation", func(t *testing.T) {
		tx := CreateTestDBWithTx(t)

		err := tx.AutoMigrate(&User{})
		require.NoError(t, err)

		// Create user in transaction
		user := User{Name: "TX User"}
		err = tx.Create(&user).Error
		require.NoError(t, err)

		// User exists in transaction
		var found User
		err = tx.First(&found, user.ID).Error
		require.NoError(t, err)
		assert.Equal(t, "TX User", found.Name)
	})

	t.Run("Multiple tests isolated", func(t *testing.T) {
		tx1 := CreateTestDBWithTx(t)
		tx2 := CreateTestDBWithTx(t)

		err := tx1.AutoMigrate(&User{})
		require.NoError(t, err)
		err = tx2.AutoMigrate(&User{})
		require.NoError(t, err)

		// Create users in separate transactions
		user1 := User{Name: "User 1"}
		user2 := User{Name: "User 2"}

		err = tx1.Create(&user1).Error
		require.NoError(t, err)
		err = tx2.Create(&user2).Error
		require.NoError(t, err)

		// Each transaction only sees its own data
		var count1, count2 int64
		tx1.Model(&User{}).Count(&count1)
		tx2.Model(&User{}).Count(&count2)

		assert.Equal(t, int64(1), count1)
		assert.Equal(t, int64(1), count2)
	})
}
