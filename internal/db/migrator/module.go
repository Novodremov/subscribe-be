package migrator

import (
	"context"

	"go.uber.org/fx"
)

var Module = fx.Module("migrator", fx.Provide(NewMigrator), fx.Invoke(func(lc fx.Lifecycle, migrator *Migrator) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return migrator.Run()
		},
	})
}))
