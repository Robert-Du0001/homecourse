package bootstrap

import (
	"github.com/goravel/framework/contracts/validation"

	"homecourse/app/rules"
)

func Rules() []validation.Rule {
	return []validation.Rule{
		&rules.Exists{},
	}
}
