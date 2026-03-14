package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"homecourse/app/facades"
)

type M20260209200954CreateCategoriesTable struct{}

// Signature The unique signature for the migration.
func (r *M20260209200954CreateCategoriesTable) Signature() string {
	return "20260209200954_create_categories_table"
}

// Up Run the migrations.
func (r *M20260209200954CreateCategoriesTable) Up() error {
	if !facades.Schema().HasTable("categories") {
		return facades.Schema().Create("categories", func(table schema.Blueprint) {
			table.ID()                       // 自增 ID
			table.String("name")             // 分类名称 (如: 计算机科学)
			table.Boolean("is_default")      // 是否为默认分类（文件会默认扫描到此分类中）
			table.Integer("sort").Default(0) // 排序
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20260209200954CreateCategoriesTable) Down() error {
	return facades.Schema().DropIfExists("categories")
}
