package usecase

import "context"

// TransferUseCase は口座間送金を実行する。トランザクション境界は RunInTx で1単位。
type TransferUseCase struct {
	tx  TxManager
	repo AccountRepository
}

// NewTransferUseCase は TransferUseCase を構築する
func NewTransferUseCase(tx TxManager, repo AccountRepository) *TransferUseCase {
	return &TransferUseCase{tx: tx, repo: repo}
}

// Execute は fromID から toID へ amount を送金する。残高不足の場合は domain.ErrInsufficientBalance を返す。
func (uc *TransferUseCase) Execute(ctx context.Context, fromID, toID string, amount int64) error {
	return uc.tx.RunInTx(ctx, func(ctx context.Context) error {
		from, err := uc.repo.GetByID(ctx, fromID)
		if err != nil {
			return err
		}
		to, err := uc.repo.GetByID(ctx, toID)
		if err != nil {
			return err
		}
		if err := from.Debit(ctx, amount); err != nil {
			return err // 不変条件違反時は ErrInsufficientBalance
		}
		if err := to.Credit(ctx, amount); err != nil {
			return err
		}
		if err := uc.repo.Save(ctx, from); err != nil {
			return err
		}
		if err := uc.repo.Save(ctx, to); err != nil {
			return err
		}
		return nil
	})
}
