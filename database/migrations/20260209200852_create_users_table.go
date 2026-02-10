package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"homecourse/app/facades"
)

type M20260209200852CreateUsersTable struct{}

// Signature The unique signature for the migration.
func (r *M20260209200852CreateUsersTable) Signature() string {
	return "20260209200852_create_users_table"
}

// Up Run the migrations.
func (r *M20260209200852CreateUsersTable) Up() error {
	if !facades.Schema().HasTable("users") {
		return facades.Schema().Create("users", func(table schema.Blueprint) {
			table.ID("id")
			table.String("name")
			table.String("password")
			table.TinyInteger("role") // 角色
			table.TimestampsTz()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20260209200852CreateUsersTable) Down() error {
	return facades.Schema().DropIfExists("users")
}
