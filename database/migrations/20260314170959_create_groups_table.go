package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"homecourse/app/facades"
)

type M20260314170959CreateGroupsTable struct{}

// Signature The unique signature for the migration.
func (r *M20260314170959CreateGroupsTable) Signature() string {
	return "20260314170959_create_groups_table"
}

// Up Run the migrations.
func (r *M20260314170959CreateGroupsTable) Up() error {
	if !facades.Schema().HasTable("groups") {
		return facades.Schema().Create("groups", func(table schema.Blueprint) {
			table.ID()                               // 自增 ID
			table.BigInteger("course_id").Default(0) // 所属课程 ID
			table.String("name")                     // 分组名称
			table.Boolean("is_default")              // 是否为默认分组（文件会默认扫描到此分组中）
			table.Integer("sort").Default(0)         // 排序
			table.TimestampsTz()

			table.Foreign("course_id").References("id").On("courses").CascadeOnDelete()

			table.Index("course_id")
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20260314170959CreateGroupsTable) Down() error {
	return facades.Schema().DropIfExists("groups")
}
