package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20260108073508CreateCoursesTable struct{}

// Signature The unique signature for the migration.
func (r *M20260108073508CreateCoursesTable) Signature() string {
	return "20260108073508_create_courses_table"
}

// Up Run the migrations.
func (r *M20260108073508CreateCoursesTable) Up() error {
	if !facades.Schema().HasTable("courses") {
		return facades.Schema().Create("courses", func(table schema.Blueprint) {
			table.ID()                                     // 自增 ID
			table.BigInteger("user_id")                    // 创建者 ID
			table.BigInteger("category_id").Default(0)     // 所属分类 ID
			table.String("title")                          // 课程标题
			table.Text("description").Nullable()           // 课程简介
			table.String("cover_path").Nullable()          // 封面图路径或 URL
			table.UnsignedTinyInteger("status").Default(0) // 状态 (0: 未开始, 1: 学习中, 2: 已完结)
			table.TimestampsTz()

			table.Index("category_id")
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20260108073508CreateCoursesTable) Down() error {
	return facades.Schema().DropIfExists("courses")
}
