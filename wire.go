//go:build wireinject
// +build wireinject

package main

import (
	"KnowEase/controllers"
	"KnowEase/dao"
	"KnowEase/routes"
	"KnowEase/services"
	"KnowEase/middleware"
	"github.com/google/wire"
)

func InitializeApp() *routes.APP {
	wire.Build(
		dao.ProviderSet,
		services.ProviderSet,
		controllers.ProviderSet,
		routes.ProviderSet,
		middleware.NewMiddleWare,
		routes.NewApp,
		ProvideDBConnectionString,
	)
	return nil
}
