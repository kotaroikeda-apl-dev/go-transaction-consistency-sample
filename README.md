# 金融ドメインの整合性設計サンプル（Go）

口座間送金を題材に、**トランザクション境界**・**集約境界**・**不変条件**をコードで示すサンプルです。

## 設計の要点

| 観点                     | 実装                                                                                                                                   |
| ------------------------ | -------------------------------------------------------------------------------------------------------------------------------------- |
| **不変条件**             | 残高はマイナスにならない。`domain.Account#Debit` で残高不足時は `ErrInsufficientBalance` を返し、残高は変更しない。                    |
| **トランザクション境界** | UseCase 単位で `TxManager.RunInTx(ctx, fn)` 内に「取得→Debit/Credit→Save」をまとめ、1単位でコミット/ロールバックできる形にしている。   |
| **集約**                 | `Account` を集約ルートとし、残高は `Balance()` で参照のみ。変更は `Debit`/`Credit` 経由のみ（外部から `balance` を直接変更できない）。 |

## ディレクトリ構成

```
.
├── cmd/
│   └── demo/
│       └── main.go          # 送金デモのエントリポイント
├── domain/
│   ├── account.go           # 集約ルート Account（Debit/Credit/Balance）
│   ├── account_test.go      # ドメイン単体テスト（不変条件など）
│   └── errors.go            # ドメインエラー（残高不足・不正金額など）
├── usecase/
│   ├── repository.go        # AccountRepository インターフェース（ポート）
│   ├── tx.go                # TxManager（RunInTx でトランザクション境界）
│   ├── transfer.go          # 送金 UseCase
│   └── transfer_test.go     # UseCase テスト（送金成功・残高不足でロールバック的挙動）
├── infra/
│   └── inmemory/
│       └── account_repository.go  # InMemory リポジトリ（DB は未使用）
├── go.mod
└── README.md
```

## 実行方法

### デモの実行

```bash
go run ./cmd/demo
```

送金成功と、残高不足時の拒否（不変条件の維持）が標準出力で確認できます。

### テスト

```bash
go test ./...
```

- `domain`: 残高不変条件（Debit で残高超過時エラー）・不正金額の単体テスト
- `usecase`: 送金 UseCase のテスト（正常送金と残高不足で失敗するケース）

## 今後の拡張

- 永続化層（DB）を入れる場合は `infra` に別パッケージを追加し、`TxManager` を DB トランザクションで実装する。
- リポジトリインターフェースは `usecase` にあり、`infra` が実装する形のまま利用できる。
