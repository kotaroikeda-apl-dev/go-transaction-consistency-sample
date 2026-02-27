package domain

import (
	"context"
	"testing"
)

func TestAccount_Debit_Credit_Invariant(t *testing.T) {
	ctx := context.Background()
	acc := NewAccount("acc-1", 100)

	// 正常: 引き落とし後は残高が減る
	if err := acc.Debit(ctx, 30); err != nil {
		t.Fatalf("Debit: %v", err)
	}
	if got := acc.Balance(); got != 70 {
		t.Errorf("Balance after Debit 30: got %d, want 70", got)
	}

	// 不変条件: 残高を超える引き落としは拒否される
	if err := acc.Debit(ctx, 100); err != ErrInsufficientBalance {
		t.Errorf("Debit over balance: got err %v, want ErrInsufficientBalance", err)
	}
	if got := acc.Balance(); got != 70 {
		t.Errorf("Balance must be unchanged: got %d, want 70", got)
	}

	// 入金
	if err := acc.Credit(ctx, 50); err != nil {
		t.Fatalf("Credit: %v", err)
	}
	if got := acc.Balance(); got != 120 {
		t.Errorf("Balance after Credit 50: got %d, want 120", got)
	}
}

func TestAccount_InvalidAmount(t *testing.T) {
	ctx := context.Background()
	acc := NewAccount("acc-1", 100)

	if err := acc.Debit(ctx, 0); err != ErrInvalidAmount {
		t.Errorf("Debit(0): got %v, want ErrInvalidAmount", err)
	}
	if err := acc.Credit(ctx, -1); err != ErrInvalidAmount {
		t.Errorf("Credit(-1): got %v, want ErrInvalidAmount", err)
	}
}
