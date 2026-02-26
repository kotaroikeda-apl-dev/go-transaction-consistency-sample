package usecase

import "context"

// TxManager はトランザクション境界を表す。UseCase は RunInTx 内で1単位として実行する。
type TxManager interface {
	RunInTx(ctx context.Context, fn func(ctx context.Context) error) error
}
