package transaction

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestTransactionContext(t *testing.T) {
	t.Parallel()

	// Setup in-memory database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	t.Run("GetTx returns nil when no transaction in context", func(t *testing.T) {
		ctx := context.Background()
		tx := GetTx(ctx)
		assert.Nil(t, tx)
	})

	t.Run("SetTx and GetTx work together", func(t *testing.T) {
		ctx := context.Background()

		// Start a transaction
		tx := db.Begin()

		// Set it in context
		ctx = SetTx(ctx, tx)

		// Retrieve it
		retrievedTx := GetTx(ctx)
		assert.NotNil(t, retrievedTx)
		assert.Equal(t, tx, retrievedTx)

		// Clean up
		tx.Rollback()
	})

	t.Run("GetTxOrDefault uses transaction when available", func(t *testing.T) {
		ctx := context.Background()

		dbFunc := GetTxOrDefault(db)

		// Without transaction, should use default
		result1 := dbFunc(ctx)
		assert.NotNil(t, result1)

		// With transaction, should use transaction
		tx := db.Begin()
		ctx = SetTx(ctx, tx)
		result2 := dbFunc(ctx)
		assert.NotNil(t, result2)

		tx.Rollback()
	})

	t.Run("Fix always uses provided database", func(t *testing.T) {
		ctx := context.Background()

		dbFunc := Fix(db)

		// Should always return the fixed database, even with transaction in context
		tx := db.Begin()
		ctx = SetTx(ctx, tx)

		result := dbFunc(ctx)
		assert.NotNil(t, result)

		tx.Rollback()
	})

	t.Run("SelectForUpdate context flag", func(t *testing.T) {
		ctx := context.Background()

		// Initially false
		assert.False(t, IsSelectForUpdate(ctx))

		// Set to true
		ctx = SelectForUpdate(ctx)
		assert.True(t, IsSelectForUpdate(ctx))

		// Transaction with SELECT FOR UPDATE
		tx := db.Begin()
		ctx = SetTx(ctx, tx)

		retrievedTx := GetTx(ctx)
		assert.NotNil(t, retrievedTx)

		tx.Rollback()
	})
}

// Example usage in a repository
type User struct {
	ID      uint `gorm:"primaryKey"`
	Name    string
	Balance int64
}

type UserRepository struct {
	db func(ctx context.Context) *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: GetTxOrDefault(db), // Use transaction if available, otherwise default DB
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *User) error {
	return r.db(ctx).Create(user).Error
}

func (r *UserRepository) UpdateBalance(ctx context.Context, userID uint, newBalance int64) error {
	return r.db(ctx).Model(&User{}).Where("id = ?", userID).Update("balance", newBalance).Error
}

func (r *UserRepository) GetUser(ctx context.Context, userID uint) (*User, error) {
	var user User
	err := r.db(ctx).First(&user, userID).Error
	return &user, err
}

// Example service using transactions
type UserService struct {
	db       *gorm.DB
	userRepo *UserRepository
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db:       db,
		userRepo: NewUserRepository(db),
	}
}

func (s *UserService) TransferBalance(ctx context.Context, fromUserID, toUserID uint, amount int64) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Inject transaction into context
		ctx = SetTx(ctx, tx)

		// Get users (using transaction)
		fromUser, err := s.userRepo.GetUser(ctx, fromUserID)
		if err != nil {
			return err
		}

		toUser, err := s.userRepo.GetUser(ctx, toUserID)
		if err != nil {
			return err
		}

		// Check sufficient balance
		if fromUser.Balance < amount {
			return assert.AnError // insufficient balance
		}

		// Update balances (using transaction)
		if err := s.userRepo.UpdateBalance(ctx, fromUserID, fromUser.Balance-amount); err != nil {
			return err
		}

		return s.userRepo.UpdateBalance(ctx, toUserID, toUser.Balance+amount)
	})
}

func TestRepositoryWithTransaction(t *testing.T) {
	// Setup database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Auto migrate
	err = db.AutoMigrate(&User{})
	require.NoError(t, err)

	// Setup service and repository
	service := NewUserService(db)
	repo := service.userRepo

	t.Run("Repository works without transaction", func(t *testing.T) {
		ctx := context.Background()

		user := &User{Name: "John", Balance: 1000}
		err := repo.CreateUser(ctx, user)
		require.NoError(t, err)

		retrieved, err := repo.GetUser(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, "John", retrieved.Name)
		assert.Equal(t, int64(1000), retrieved.Balance)
	})

	t.Run("Repository works with transaction", func(t *testing.T) {
		ctx := context.Background()

		err := db.Transaction(func(tx *gorm.DB) error {
			ctx = SetTx(ctx, tx)

			user := &User{Name: "Jane", Balance: 2000}
			return repo.CreateUser(ctx, user)
		})
		require.NoError(t, err)
	})

	t.Run("Service transfer uses transaction correctly", func(t *testing.T) {
		ctx := context.Background()

		// Create test users
		user1 := &User{Name: "Alice", Balance: 1000}
		user2 := &User{Name: "Bob", Balance: 500}

		require.NoError(t, repo.CreateUser(ctx, user1))
		require.NoError(t, repo.CreateUser(ctx, user2))

		// Transfer money
		err := service.TransferBalance(ctx, user1.ID, user2.ID, 300)
		require.NoError(t, err)

		// Check final balances
		finalUser1, err := repo.GetUser(ctx, user1.ID)
		require.NoError(t, err)
		assert.Equal(t, int64(700), finalUser1.Balance)

		finalUser2, err := repo.GetUser(ctx, user2.ID)
		require.NoError(t, err)
		assert.Equal(t, int64(800), finalUser2.Balance)
	})
}
