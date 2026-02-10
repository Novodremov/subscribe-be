package repo

import (

	"go.uber.org/fx"
)

var Module = fx.Module("repository", fx.Provide(
	NewSubscriptionRepo,
))
