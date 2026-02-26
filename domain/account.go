package domain

import (
	"context"
	"sync"
)

// Account は集約ルート。残高はこの型のメソッド経由でのみ変更する（不変条件: 残高 >= 0）。
type Account struct {
	mu      sync.Mutex
	id      string
	balance int64
}

// NewAccount は永続化層から再構成するとき、または新規作成時に使う
func NewAccount(id string, balance int64) *Account {
	return &Account{id: id, balance: balance}
}

// ID は口座IDを返す
func (a *Account) ID() string {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.id
}

// Balance は現在残高を返す（読み取り専用）
func (a *Account) Balance() int64 {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.balance
}

// Debit は口座から指定額を引き落とす。残高がマイナスになる場合は ErrInsufficientBalance を返す（不変条件の保証）。
func (a *Account) Debit(_ context.Context, amount int64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.balance < amount {
		return ErrInsufficientBalance
	}
	a.balance -= amount
	return nil
}

// Credit は口座に指定額を入金する。
func (a *Account) Credit(_ context.Context, amount int64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.balance += amount
	return nil
}
