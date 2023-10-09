package sub_process

import (
	"github.com/google/wire"
	"gitlab.yoyiit.com/youyi/app-bank/internal/repo"
	"gitlab.yoyiit.com/youyi/app-base/kitex_gen/api/base"
	"gitlab.yoyiit.com/youyi/app-oa/kitex_gen/api/oa"
)

var ProviderSet = wire.NewSet(NewPaymentReceiptSubProcess)

func NewPaymentReceiptSubProcess(paymentReceiptRepo repo.PaymentReceiptRepo, oaClient oa.Client, baseClient base.Client,
	paymentReceiptApplicationCustomFieldRepo repo.PaymentReceiptApplicationCustomFieldRepo) *PaymentReceiptSubProcess {
	return &PaymentReceiptSubProcess{paymentReceiptRepo, oaClient, baseClient,
		paymentReceiptApplicationCustomFieldRepo}
}
