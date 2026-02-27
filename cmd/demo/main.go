// 口座間送金デモ。トランザクション境界（RunInTx）・集約（Account）・不変条件（残高>=0）の動作を示す。
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ikedakotarou/go-transaction-consistency-sample/domain"
	"github.com/ikedakotarou/go-transaction-consistency-sample/infra/inmemory"
	"github.com/ikedakotarou/go-transaction-consistency-sample/usecase"
)

func main() {
	ctx := context.Background()
	repo := inmemory.NewAccountRepository()

	// 初期データ: 口座 A=1000, B=200
	_ = repo.Save(ctx, domain.NewAccount("A", 1000))
	_ = repo.Save(ctx, domain.NewAccount("B", 200))

	tx := &inmemoryTx{}
	transfer := usecase.NewTransferUseCase(tx, repo)

	fmt.Println("--- 送金デモ（トランザクション境界・集約・不変条件） ---")
	printBalances(repo, ctx, "A", "B")

	// 正常送金: A -> B に 300
	if err := transfer.Execute(ctx, "A", "B", 300); err != nil {
		log.Fatalf("transfer: %v", err)
	}
	fmt.Println("A -> B に 300 送金後:")
	printBalances(repo, ctx, "A", "B")

	// 不変条件違反: 残高を超える送金は拒否され、残高は変化しない
	if err := transfer.Execute(ctx, "A", "B", 1000); err != nil {
		fmt.Printf("残高不足で送金拒否（期待通り）: %v\n", err)
	}
	fmt.Println("拒否後の残高（変化なし）:")
	printBalances(repo, ctx, "A", "B")
}

func printBalances(repo *inmemory.AccountRepository, ctx context.Context, ids ...string) {
	for _, id := range ids {
		acc, err := repo.GetByID(ctx, id)
		if err != nil {
			fmt.Printf("  %s: error %v\n", id, err)
			continue
		}
		fmt.Printf("  %s: %d\n", id, acc.Balance())
	}
}

type inmemoryTx struct{}

func (t *inmemoryTx) RunInTx(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}
