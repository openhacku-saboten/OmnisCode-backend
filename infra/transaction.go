package infra

import (
	"context"

	"github.com/go-gorp/gorp"
)

// ref: https://qiita.com/miya-masa/items/316256924a1f0d7374bb
var txKey = struct{}{}

// TransactionDAO はTransactionに関するDataAccessObjectです
// gorp.dbMapで利用するメソッドをここにインタフェースとして定義することで、
// infraではgorpに依存しないような設計となっています
type TransactionDAO interface {
	// ref: https://pkg.go.dev/github.com/go-gorp/gorp#DbMap.Delete
	Delete(list ...interface{}) (int64, error)
}

// getTx はcontextからトランザクションを取得する
func getTx(ctx context.Context) (TransactionDAO, bool) {
	tx, ok := ctx.Value(&txKey).(*gorp.Transaction)
	return tx, ok
}
