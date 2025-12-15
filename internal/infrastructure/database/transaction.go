package database

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// TransactionManager handles database transactions
type TransactionManager struct {
	db *gorm.DB
}

// NewTransactionManager creates a new transaction manager
func NewTransactionManager(db *gorm.DB) *TransactionManager {
	return &TransactionManager{db: db}
}

// WithTransaction executes a function within a database transaction
func (tm *TransactionManager) WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	tx := tm.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r) // re-throw panic after rollback
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback().Error; rbErr != nil {
			return fmt.Errorf("failed to rollback transaction: %w (original error: %v)", rbErr, err)
		}
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// TransactionalRepository interface that repositories should implement
type TransactionalRepository interface {
	// WithTx returns a repository instance that uses the given transaction
	WithTx(tx *gorm.DB) interface{}
}

// GetDB returns the underlying database instance
func (tm *TransactionManager) GetDB() *gorm.DB {
	return tm.db
}

// InTransaction checks if the context is already in a transaction
func InTransaction(db *gorm.DB) bool {
	return db.Statement.DB.Statement.DB != nil
}
