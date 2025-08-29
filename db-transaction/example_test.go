package transaction

import (
	"context"
	"fmt"
	"testing"

	dbtesting "db-testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// Account represents a user account for the banking example
type Account struct {
	ID      uint   `gorm:"primaryKey"`
	Name    string `gorm:"not null"`
	Balance int64  `gorm:"default:0"`
}

// AccountRepository handles account data operations
type AccountRepository struct {
	db func(ctx context.Context) *gorm.DB
}

// NewAccountRepository creates a new account repository
func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{
		db: GetTxOrDefault(db),
	}
}

func (r *AccountRepository) CreateAccount(ctx context.Context, account *Account) error {
	return r.db(ctx).Create(account).Error
}

func (r *AccountRepository) GetAccount(ctx context.Context, id uint) (*Account, error) {
	var account Account
	err := r.db(ctx).First(&account, id).Error
	return &account, err
}

func (r *AccountRepository) UpdateBalance(ctx context.Context, id uint, newBalance int64) error {
	return r.db(ctx).Model(&Account{}).Where("id = ?", id).Update("balance", newBalance).Error
}

// BankingService handles business logic for banking operations
type BankingService struct {
	db      *gorm.DB
	accRepo *AccountRepository
}

func NewBankingService(db *gorm.DB) *BankingService {
	return &BankingService{
		db:      db,
		accRepo: NewAccountRepository(db),
	}
}

// TransferMoney transfers money between accounts atomically
func (s *BankingService) TransferMoney(ctx context.Context, fromID, toID uint, amount int64) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Inject transaction into context - this is the key pattern!
		ctx = SetTx(ctx, tx)

		// All repository operations will now use the transaction automatically
		fromAccount, err := s.accRepo.GetAccount(ctx, fromID)
		if err != nil {
			return fmt.Errorf("failed to get source account: %w", err)
		}

		toAccount, err := s.accRepo.GetAccount(ctx, toID)
		if err != nil {
			return fmt.Errorf("failed to get destination account: %w", err)
		}

		// Business logic validation
		if fromAccount.Balance < amount {
			return fmt.Errorf("insufficient balance: has %d, needs %d", fromAccount.Balance, amount)
		}

		// Update balances
		if err := s.accRepo.UpdateBalance(ctx, fromID, fromAccount.Balance-amount); err != nil {
			return fmt.Errorf("failed to debit source account: %w", err)
		}

		if err := s.accRepo.UpdateBalance(ctx, toID, toAccount.Balance+amount); err != nil {
			return fmt.Errorf("failed to credit destination account: %w", err)
		}

		return nil
	})
}

// CreateAccountWithInitialDeposit creates account and sets initial balance atomically
func (s *BankingService) CreateAccountWithInitialDeposit(ctx context.Context, name string, initialBalance int64) (*Account, error) {
	var account *Account

	err := s.db.Transaction(func(tx *gorm.DB) error {
		ctx = SetTx(ctx, tx)

		account = &Account{
			Name:    name,
			Balance: initialBalance,
		}

		return s.accRepo.CreateAccount(ctx, account)
	})

	return account, err
}

// TestBankingTransactionExample demonstrates the complete banking transaction pattern
func TestBankingTransactionExample(t *testing.T) {
	// Setup PostgreSQL database
	db := dbtesting.CreateTestDB(t, dbtesting.EnvTest, dbtesting.DBDebugOff)

	// Auto migrate
	require.NoError(t, db.AutoMigrate(&Account{}))

	// Create banking service
	bankingService := NewBankingService(db)
	ctx := context.Background()

	t.Run("Complete Banking Transaction Demo", func(t *testing.T) {
		// Create accounts with initial deposits (using transactions)
		alice, err := bankingService.CreateAccountWithInitialDeposit(ctx, "Alice", 1000)
		require.NoError(t, err)
		require.Equal(t, "Alice", alice.Name)
		require.Equal(t, int64(1000), alice.Balance)

		bob, err := bankingService.CreateAccountWithInitialDeposit(ctx, "Bob", 500)
		require.NoError(t, err)
		require.Equal(t, "Bob", bob.Name)
		require.Equal(t, int64(500), bob.Balance)

		// Transfer money (using transactions)
		err = bankingService.TransferMoney(ctx, alice.ID, bob.ID, 300)
		require.NoError(t, err)

		// Check final balances
		finalAlice, err := bankingService.accRepo.GetAccount(ctx, alice.ID)
		require.NoError(t, err)
		require.Equal(t, int64(700), finalAlice.Balance)

		finalBob, err := bankingService.accRepo.GetAccount(ctx, bob.ID)
		require.NoError(t, err)
		require.Equal(t, int64(800), finalBob.Balance)
	})

	t.Run("Transaction Rollback on Insufficient Funds", func(t *testing.T) {
		// Create test accounts
		charlie, err := bankingService.CreateAccountWithInitialDeposit(ctx, "Charlie", 100)
		require.NoError(t, err)

		dave, err := bankingService.CreateAccountWithInitialDeposit(ctx, "Dave", 200)
		require.NoError(t, err)

		// Record initial balances
		initialCharlie := charlie.Balance
		initialDave := dave.Balance

		// Attempt transfer with insufficient funds (should fail and rollback)
		err = bankingService.TransferMoney(ctx, charlie.ID, dave.ID, 1000)
		require.Error(t, err)
		require.Contains(t, err.Error(), "insufficient balance")

		// Verify balances didn't change (transaction rolled back)
		finalCharlie, err := bankingService.accRepo.GetAccount(ctx, charlie.ID)
		require.NoError(t, err)
		require.Equal(t, initialCharlie, finalCharlie.Balance)

		finalDave, err := bankingService.accRepo.GetAccount(ctx, dave.ID)
		require.NoError(t, err)
		require.Equal(t, initialDave, finalDave.Balance)
	})

	t.Run("Repository Works Without Transaction", func(t *testing.T) {
		// Repository methods work fine without transactions too
		eve := &Account{Name: "Eve", Balance: 750}
		err := bankingService.accRepo.CreateAccount(ctx, eve)
		require.NoError(t, err)

		retrieved, err := bankingService.accRepo.GetAccount(ctx, eve.ID)
		require.NoError(t, err)
		require.Equal(t, "Eve", retrieved.Name)
		require.Equal(t, int64(750), retrieved.Balance)
	})
}
