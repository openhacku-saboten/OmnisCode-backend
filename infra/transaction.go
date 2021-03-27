package infra

import (
	"context"

	"github.com/go-gorp/gorp"
)

// ref: https://qiita.com/miya-masa/items/316256924a1f0d7374bb#%E8%A7%A3%E6%B1%BA%E6%A1%88%EF%BC%92-%E3%82%B3%E3%83%B3%E3%83%86%E3%82%AD%E3%82%B9%E3%83%88%E3%81%AB%E3%83%88%E3%83%A9%E3%83%B3%E3%82%B6%E3%82%AF%E3%82%B7%E3%83%A7%E3%83%B3%E3%82%AA%E3%83%96%E3%82%B8%E3%82%A7%E3%82%AF%E3%83%88%E3%82%92%E3%82%BB%E3%83%83%E3%83%88%E3%81%99%E3%82%8B
var txKey = struct{}{}

// TransactionDAO はTransactionに関するDataAccessObjectです
type TransactionDAO interface {
}

// getTx はcontextからトランザクションを取得する
func getTx(ctx context.Context) (TransactionDAO, bool) {
	tx, ok := ctx.Value(&txKey).(*gorp.Transaction)
	return tx, ok
}
