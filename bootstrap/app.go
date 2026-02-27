package bootstrap

import (
	contractsfoundation "github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/foundation"

	"homecourse/config"
	"homecourse/routes"
)

func Boot() contractsfoundation.Application {
	return foundation.Setup().
		WithCommands(Commands).
		WithMigrations(Migrations).
		WithRouting(func() {
			routes.Web()
			routes.Api()
		}).
		WithProviders(Providers).
		WithConfig(config.Boot).
		Create()
}
