package service

import "go.uber.org/fx"

var Module = fx.Module("svc",
	fx.Provide(
		NewSubscriptionService,
	),
)
