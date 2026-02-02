package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20260108073203CreateCategoriesTable struct{}

// Signature The unique signature for the migration.
func (r *M20260108073203CreateCategoriesTable) Signature() string {
	return "20260108073203_create_categories_table"
}

// Up Run the migrations.
func (r *M20260108073203CreateCategoriesTable) Up() error {
	if !facades.Schema().HasTable("categories") {
		return facades.Schema().Create("categories", func(table schema.Blueprint) {
			table.ID()                               // 自增 ID
			table.BigInteger("user_id")              // 创建者 ID
			table.String("name")                     // 分类名称 (如: 计算机科学)
			table.BigInteger("parent_id").Default(0) // 父分类 ID (默认为 0，方便做二级分类)
			table.Integer("sort").Default(0)         // 排序
			table.TimestampsTz()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20260108073203CreateCategoriesTable) Down() error {
	return facades.Schema().DropIfExists("categories")
}
