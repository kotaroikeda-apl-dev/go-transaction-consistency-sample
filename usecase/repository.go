package usecase

import (
	"context"

	"github.com/ikedakotarou/go-transaction-consistency-sample/domain"
)

// AccountRepository は Account 集約の永続化（読み取り・保存）を抽象化する
type AccountRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Account, error)
	Save(ctx context.Context, account *domain.Account) error
}
