package main

import (
	"github.com/cloudwego/kitex/server"
	"gitlab.yoyiit.com/youyi/app-bank/internal/service"
	"gitlab.yoyiit.com/youyi/app-bank/kitex_gen/api"
	"gitlab.yoyiit.com/youyi/app-bank/kitex_gen/api/bank"
	"gitlab.yoyiit.com/youyi/go-common/enum"
	"gitlab.yoyiit.com/youyi/go-core/launch"
)

func main() {
	logger, closer := launch.InitPremise()
	defer logger.Sync()
	defer closer.Close()

	enum.Init()
	server, err := initServer()
	if err != nil {
		panic(err)
	}
	go launch.RunServer(server)
	launch.InitHttpServer()
}

func newBankImpl(bankService service.BankService, paymentReceiptService service.PaymentReceiptService) api.Bank {
	return &BankImpl{bankService, paymentReceiptService}
}

func newServer(handler api.Bank) server.Server {
	return bank.NewServer(handler, launch.RpcServerOptions()...)
}
