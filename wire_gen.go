// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/cloudwego/kitex/server"
	"gitlab.yoyiit.com/youyi/app-bank/internal/repo"
	"gitlab.yoyiit.com/youyi/app-bank/internal/sdk"
	"gitlab.yoyiit.com/youyi/app-bank/internal/service"
	"gitlab.yoyiit.com/youyi/app-bank/internal/sub_process"
	"gitlab.yoyiit.com/youyi/app-base/process"
	"gitlab.yoyiit.com/youyi/app-base/rpc_base"
	"gitlab.yoyiit.com/youyi/app-dingtalk/rpc_dingtalk"
	"gitlab.yoyiit.com/youyi/app-finance/rpc_finance"
	"gitlab.yoyiit.com/youyi/app-invoice/rpc_invoice"
	"gitlab.yoyiit.com/youyi/app-oa/rpc_oa"
	"gitlab.yoyiit.com/youyi/app-soms/rpc_soms"
	"gitlab.yoyiit.com/youyi/go-core/store"
	"gitlab.yoyiit.com/youyi/go-core/trace"
)

// Injectors from wire.go:

func initServer() (server.Server, error) {
	client := rpc_base.NewBaseClient()
	guilinBankSDK := sdk.NewGuilinBankSDK()
	spdBankSDK := sdk.NewSPDBankSDK()
	pinganBankSDK := sdk.NewPinganBankSDK()
	loggerInterface := trace.NewGormLogger()
	db, err := store.NewReadWriteSeparationDB(loggerInterface)
	if err != nil {
		return nil, err
	}
	bankTransferReceiptRepo := repo.NewBankTransferReceiptRepo(db)
	kafkaProducer := store.NewKafkaProducer()
	dingtalkClient := rpc_dingtalk.NewDingtalkClient()
	bankTransactionDetailRepo := repo.NewBankTransactionDetailRepo(db)
	bankTransactionDetailProcessInstanceRepo := repo.NewBankTransactionDetailExternalRepo(db)
	ossConfig := store.NewOSSConfig()
	bankCodeRepo := repo.NewBankCodeRepo(db)
	bankBusinessPayrollRepo := repo.NewBankBusinessPayrollRepo(db)
	bankBusinessPayrollDetailRepo := repo.NewBankBusinessPayrollDetailRepo(db)
	oaClient := rpc_oa.NewOAClient()
	paymentReceiptRepo := repo.NewPaymentReceiptRepo(db)
	pdfToImageService := service.NewPdfToImageService(ossConfig)
	financeClient := rpc_finance.NewFinanceClient()
	icbcBankSDK := sdk.NewIcbcBankSDK()
	redisClient := store.NewRedisClient()
	minShengSDK := sdk.NewMinShengBankSDK()
	bankService := service.NewBankService(client, guilinBankSDK, spdBankSDK, pinganBankSDK, bankTransferReceiptRepo, kafkaProducer, dingtalkClient, bankTransactionDetailRepo, bankTransactionDetailProcessInstanceRepo, ossConfig, bankCodeRepo, bankBusinessPayrollRepo, bankBusinessPayrollDetailRepo, oaClient, paymentReceiptRepo, pdfToImageService, financeClient, icbcBankSDK, redisClient, minShengSDK)
	paymentReceiptApplicationCustomFieldRepo := repo.NewPaymentReceiptApplicationCustomFieldRepo(db)
	invoiceClient := rpc_invoice.NewInvoiceClient()
	paymentReceiptSubProcess := sub_process.NewPaymentReceiptSubProcess(paymentReceiptRepo, oaClient, client, paymentReceiptApplicationCustomFieldRepo, invoiceClient)
	processAuthRepo := process.NewProcessAuthRepo(db)
	somsClient := rpc_soms.NewSomsClient()
	paymentReceiptService := service.NewPaymentReceiptService(paymentReceiptSubProcess, client, paymentReceiptRepo, bankCodeRepo, guilinBankSDK, spdBankSDK, pinganBankSDK, oaClient, dingtalkClient, processAuthRepo, somsClient, paymentReceiptApplicationCustomFieldRepo, invoiceClient, minShengSDK)
	bank := newBankImpl(bankService, paymentReceiptService)
	serverServer := newServer(bank)
	return serverServer, nil
}
