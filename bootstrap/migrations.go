package bootstrap

import (
	"github.com/goravel/framework/contracts/database/schema"

	"homecourse/database/migrations"
)

func Migrations() []schema.Migration {
	return []schema.Migration{
		&migrations.M20260209200852CreateUsersTable{},
		&migrations.M20260209200954CreateCategoriesTable{},
		&migrations.M20260209201035CreateCoursesTable{},
		&migrations.M20260209201058CreateEpisodesTable{},
	}
}
