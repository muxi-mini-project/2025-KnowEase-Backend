package dao

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	ProvideEmailDao,
	ProvideLikeDao,
	ProvidePostDao,
	ProvideUserDao,
	NewDB,
)
