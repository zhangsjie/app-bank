package service

import (
	"github.com/google/wire"
	"gitlab.yoyiit.com/youyi/app-bank/internal/repo"
	"gitlab.yoyiit.com/youyi/app-bank/internal/sdk"
	"gitlab.yoyiit.com/youyi/app-bank/internal/sub_process"
	"gitlab.yoyiit.com/youyi/app-base/kitex_gen/api/base"
	"gitlab.yoyiit.com/youyi/app-base/process"
	"gitlab.yoyiit.com/youyi/app-base/rpc_base"
	"gitlab.yoyiit.com/youyi/app-dingtalk/kitex_gen/api/dingtalk"
	"gitlab.yoyiit.com/youyi/app-dingtalk/rpc_dingtalk"
	"gitlab.yoyiit.com/youyi/app-finance/kitex_gen/api/finance"
	"gitlab.yoyiit.com/youyi/app-finance/rpc_finance"
	"gitlab.yoyiit.com/youyi/app-oa/kitex_gen/api/oa"
	"gitlab.yoyiit.com/youyi/app-oa/rpc_oa"
	"gitlab.yoyiit.com/youyi/app-soms/kitex_gen/api/soms"
	"gitlab.yoyiit.com/youyi/app-soms/rpc_soms"
	"gitlab.yoyiit.com/youyi/go-core/store"
)

var ProviderSet = wire.NewSet(NewBankService, store.NewKafkaProducer, rpc_dingtalk.NewDingtalkClient, store.NewOSSConfig,
	rpc_base.NewBaseClient, NewPaymentReceiptService, rpc_oa.NewOAClient, NewPdfToImageService, process.NewProcessAuthRepo,
	rpc_soms.NewSomsClient, rpc_finance.NewFinanceClient)

func NewBankService(baseClient base.Client, guilinBankSDK sdk.GuilinBankSDK, spdBankSDK sdk.SPDBankSDK, pinganBankSDK sdk.PinganBankSDK,
	bankTransferReceiptRepo repo.BankTransferReceiptRepo, kafkaProducer *store.KafkaProducer, dingtalkClient dingtalk.Client,
	bankTransactionDetailRepo repo.BankTransactionDetailRepo, bankTransactionDetailProcessInstanceRepo repo.BankTransactionDetailProcessInstanceRepo,
	ossConfig *store.OSSConfig, bankCode repo.BankCodeRepo, businessPayrollRepo repo.BankBusinessPayrollRepo,
	businessPayrollDetailRepo repo.BankBusinessPayrollDetailRepo, oaClient oa.Client, paymentReceiptRepo repo.PaymentReceiptRepo, pdfToImageService PdfToImageService, financeClient finance.Client) BankService {
	return &bankService{baseClient, guilinBankSDK, spdBankSDK, pinganBankSDK,
		bankTransferReceiptRepo, kafkaProducer, dingtalkClient,
		bankTransactionDetailRepo, bankTransactionDetailProcessInstanceRepo,
		ossConfig, bankCode, businessPayrollRepo,
		businessPayrollDetailRepo, oaClient, paymentReceiptRepo, pdfToImageService, financeClient}
}

func NewPaymentReceiptService(paymentReceiptSubProcess *sub_process.PaymentReceiptSubProcess, baseClient base.Client,
	paymentReceiptRepo repo.PaymentReceiptRepo, bankCodeRepo repo.BankCodeRepo, guilinBankSDK sdk.GuilinBankSDK,
	spdBankSDK sdk.SPDBankSDK, pinganBankSDK sdk.PinganBankSDK, oaClient oa.Client, dingtalkClient dingtalk.Client,
	processAuthRepo process.ProcessAuthRepo, somsClient soms.Client, paymentReceiptApplicationCustomFieldRepo repo.PaymentReceiptApplicationCustomFieldRepo) PaymentReceiptService {
	return &paymentReceiptService{
		process.Process{
			SubProcess:      paymentReceiptSubProcess,
			BaseClient:      baseClient,
			ProcessAuthRepo: processAuthRepo,
		},
		paymentReceiptRepo, baseClient, bankCodeRepo, guilinBankSDK,
		spdBankSDK, pinganBankSDK, oaClient, dingtalkClient, somsClient,
		paymentReceiptApplicationCustomFieldRepo,
	}
}

func NewPdfToImageService(ossConfig *store.OSSConfig) PdfToImageService {
	return &pdfToImageService{
		ossConfig: ossConfig,
	}
}
