package bootstrap

import (
	"github.com/goravel/framework/contracts/console"

	"homecourse/app/console/commands"
)

func Commands() []console.Command {
	return []console.Command{
		&commands.Init{},
	}
}
