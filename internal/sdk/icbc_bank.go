package sdk

import (
	"context"
	"gitlab.yoyiit.com/youyi/app-bank/internal/sdk/stru"
)

type IcbcBankSDK interface {
	ListTransactionDetail(ctx context.Context, account string, beginDate string, endDate string) ([]stru.IcbcAccDetailItem, error)
}

type icbcBankSDK struct{}

func (i *icbcBankSDK) ListTransactionDetail(ctx context.Context, account string, beginDate string, endDate string) ([]stru.IcbcAccDetailItem, error) {

	return nil, nil
}
