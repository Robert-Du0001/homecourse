package database

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/contracts/database/seeder"

	"homecourse/database/migrations"
	"homecourse/database/seeders"
)

type Kernel struct {
}

func (kernel Kernel) Migrations() []schema.Migration {
	return []schema.Migration{
		&migrations.M20210101000001CreateUsersTable{},
		&migrations.M20210101000002CreateJobsTable{},
		&migrations.M20260108073203CreateCategoriesTable{},
		&migrations.M20260108073508CreateCoursesTable{},
		&migrations.M20260108092605CreateEpisodesTable{},
	}
}

func (kernel Kernel) Seeders() []seeder.Seeder {
	return []seeder.Seeder{
		&seeders.DatabaseSeeder{},
	}
}
