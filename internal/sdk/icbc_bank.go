package sdk

import (
	"context"
	"fmt"
	"gitlab.yoyiit.com/youyi/app-bank/internal/sdk/stru"
	"time"
)

type IcbcBankSDK interface {
	ListTransactionDetail(ctx context.Context, account string, beginDate string, endDate string) ([]stru.IcbcAccDetailItem, error)
}

type icbcBankSDK struct{}

func (i *icbcBankSDK) ListTransactionDetail(ctx context.Context, account string, beginDate string, endDate string) ([]stru.IcbcAccDetailItem, error) {
	return nil, nil
}

func main() {
	appID := "your_app_id"
	msgID := "your_msg_id"
	signType := "RSA"
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	bizContent := stru.AccDetailRequest{
		Fseqno:    "sample_fseqno",
		Account:   "sample_account",
		Currtype:  1,
		Startdate: "2023-07-01",
		Enddate:   "2023-07-31",
		Serialno:  "",
		Corpno:    "sample_corpno",
		Acccompno: "sample_acccompno",
		Agreeno:   "sample_agreeno",
	}

	request := stru.IcbcGlobalRequest{
		AppID:      appID,
		MsgID:      msgID,
		SignType:   signType,
		Timestamp:  timestamp,
		BizContent: bizContent,
	}

	signature := stru.GenerateSignature(request, "your_private_key")
	request.Sign = signature

	fmt.Println("Request:", request)
}
