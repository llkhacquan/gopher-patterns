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
	t.Run("EnvTest with default options", func(t *testing.T) {
		db := CreateTestDB(t, EnvTest)

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
	})

	t.Run("EnvTest with debug off", func(t *testing.T) {
		db := CreateTestDB(t, EnvTest, DBDebugOff)

		err := db.AutoMigrate(&User{})
		require.NoError(t, err)

		user := User{Name: "Debug Off User"}
		err = db.Create(&user).Error
		require.NoError(t, err)
		assert.NotZero(t, user.ID)
	})

	t.Run("EnvTest without transaction wrapping", func(t *testing.T) {
		db := CreateTestDB(t, EnvTest, DBNoWrapInTransaction)

		err := db.AutoMigrate(&User{})
		require.NoError(t, err)

		user := User{Name: "No Transaction User"}
		err = db.Create(&user).Error
		require.NoError(t, err)
		assert.NotZero(t, user.ID)
	})

	t.Run("EnvDev (may skip if not available)", func(t *testing.T) {
		db := CreateTestDB(t, EnvDev, DBDebugOff)
		if db == nil {
			t.Skip("EnvDev not available")
			return
		}

		err := db.AutoMigrate(&User{})
		require.NoError(t, err)

		user := User{Name: "Dev User"}
		err = db.Create(&user).Error
		require.NoError(t, err)
		assert.NotZero(t, user.ID)
	})
}

func TestBackwardsCompatibility(t *testing.T) {
	t.Run("Legacy CreateTestDBLegacy", func(t *testing.T) {
		db := CreateTestDBLegacy(t)

		err := db.AutoMigrate(&User{})
		require.NoError(t, err)

		user := User{Name: "Legacy User"}
		err = db.Create(&user).Error
		require.NoError(t, err)
		assert.NotZero(t, user.ID)
	})
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

func TestDatabaseOptions(t *testing.T) {
	t.Run("Combined options", func(t *testing.T) {
		// Test multiple options together
		db := CreateTestDB(t, EnvTest, DBDebugOff, DBNoWrapInTransaction)

		err := db.AutoMigrate(&User{})
		require.NoError(t, err)

		user := User{Name: "Multi Options User"}
		err = db.Create(&user).Error
		require.NoError(t, err)
		assert.NotZero(t, user.ID)

		// Verify user persists (no transaction rollback)
		var found User
		err = db.First(&found, user.ID).Error
		require.NoError(t, err)
		assert.Equal(t, "Multi Options User", found.Name)
	})

	t.Run("Connection caching", func(t *testing.T) {
		// Multiple calls should reuse cached connection
		db1 := CreateTestDB(t, EnvTest, DBDebugOff)
		db2 := CreateTestDB(t, EnvTest, DBDebugOff)

		// Both should work independently
		err := db1.AutoMigrate(&User{})
		require.NoError(t, err)
		err = db2.AutoMigrate(&User{})
		require.NoError(t, err)

		user1 := User{Name: "Cache User 1"}
		user2 := User{Name: "Cache User 2"}

		err = db1.Create(&user1).Error
		require.NoError(t, err)
		err = db2.Create(&user2).Error
		require.NoError(t, err)

		// Both users should be created successfully
		assert.NotZero(t, user1.ID)
		assert.NotZero(t, user2.ID)
		// Note: IDs may be the same since they're in separate isolated databases

		// Verify they can be found by name (proving isolation)
		var found1, found2 User
		err = db1.Where("name = ?", "Cache User 1").First(&found1).Error
		require.NoError(t, err)
		err = db2.Where("name = ?", "Cache User 2").First(&found2).Error
		require.NoError(t, err)

		assert.Equal(t, "Cache User 1", found1.Name)
		assert.Equal(t, "Cache User 2", found2.Name)
	})
}
