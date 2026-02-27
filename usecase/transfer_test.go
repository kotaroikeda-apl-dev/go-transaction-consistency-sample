package usecase

import (
	"context"
	"testing"

	"github.com/ikedakotarou/go-transaction-consistency-sample/domain"
	"github.com/ikedakotarou/go-transaction-consistency-sample/infra/inmemory"
)

func TestTransferUseCase_Execute(t *testing.T) {
	ctx := context.Background()
	repo := inmemory.NewAccountRepository()
	_ = repo.Save(ctx, domain.NewAccount("A", 100))
	_ = repo.Save(ctx, domain.NewAccount("B", 50))

	tx := &inMemoryTx{}
	uc := NewTransferUseCase(tx, repo)

	// 正常: A -> B に 30 送金
	if err := uc.Execute(ctx, "A", "B", 30); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	a, _ := repo.GetByID(ctx, "A")
	b, _ := repo.GetByID(ctx, "B")
	if a.Balance() != 70 || b.Balance() != 80 {
		t.Errorf("after transfer: A=%d, B=%d; want A=70, B=80", a.Balance(), b.Balance())
	}

	// 不変条件: 残高不足で送金は失敗し、残高は変化しない
	if err := uc.Execute(ctx, "A", "B", 100); err != domain.ErrInsufficientBalance {
		t.Errorf("Execute over balance: got err %v, want ErrInsufficientBalance", err)
	}
	a, _ = repo.GetByID(ctx, "A")
	b, _ = repo.GetByID(ctx, "B")
	if a.Balance() != 70 || b.Balance() != 80 {
		t.Errorf("balances must be unchanged: A=%d, B=%d", a.Balance(), b.Balance())
	}
}

// inMemoryTx は RunInTx をそのまま fn 実行で表現（InMemory ではトランザクションは不要）
type inMemoryTx struct{}

func (t *inMemoryTx) RunInTx(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}
