package rpc_bank

import (
	"github.com/cloudwego/kitex/client"
	"gitlab.yoyiit.com/youyi/app-bank/kitex_gen/api/bank"
	"gitlab.yoyiit.com/youyi/go-core/config"
	"gitlab.yoyiit.com/youyi/go-core/trace"
	"gitlab.yoyiit.com/youyi/go-core/util"
	"time"
)

func NewBankClient() bank.Client {
	client, err := bank.NewClient(
		config.GetString("app.bank.name", ""),
		util.GetClientOption(config.GetString("app.bank.port", "")),
		client.WithRPCTimeout(5*60*time.Second),
		client.WithSuite(trace.NewClientSuite()),
	)
	if err != nil {
		panic(err)
	}
	return client
}
