package inmemory

import (
	"context"
	"sync"

	"github.com/ikedakotarou/go-transaction-consistency-sample/domain"
)

// AccountRepository は Account をメモリ上に保持するリポジトリ実装
type AccountRepository struct {
	mu   sync.RWMutex
	byID map[string]*domain.Account
}

// NewAccountRepository は InMemory の AccountRepository を返す
func NewAccountRepository() *AccountRepository {
	return &AccountRepository{byID: make(map[string]*domain.Account)}
}

// GetByID は ID で口座を取得する。存在しない場合は ErrAccountNotFound を返す
func (r *AccountRepository) GetByID(ctx context.Context, id string) (*domain.Account, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	acc, ok := r.byID[id]
	if !ok {
		return nil, domain.ErrAccountNotFound
	}
	return acc, nil
}

// Save は口座の現在状態を保存する
func (r *AccountRepository) Save(ctx context.Context, account *domain.Account) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if account == nil {
		return nil
	}
	r.byID[account.ID()] = account
	return nil
}
