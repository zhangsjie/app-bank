//go:build wireinject
// +build wireinject

package main

import (
	"github.com/cloudwego/kitex/server"
	"github.com/google/wire"
	"gitlab.yoyiit.com/youyi/app-bank/internal/repo"
	"gitlab.yoyiit.com/youyi/app-bank/internal/sdk"
	"gitlab.yoyiit.com/youyi/app-bank/internal/service"
	"gitlab.yoyiit.com/youyi/app-bank/internal/sub_process"
)

func initServer() (server.Server, error) {
	panic(wire.Build(service.ProviderSet, sub_process.ProviderSet, sdk.ProviderSet, repo.ProviderSet, newBankImpl, newServer))
}
