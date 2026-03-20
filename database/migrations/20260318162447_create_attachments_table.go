package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"homecourse/app/facades"
)

type M20260318162447CreateAttachmentsTable struct{}

// Signature The unique signature for the migration.
func (r *M20260318162447CreateAttachmentsTable) Signature() string {
	return "20260318162447_create_attachments_table"
}

// Up Run the migrations.
func (r *M20260318162447CreateAttachmentsTable) Up() error {
	if !facades.Schema().HasTable("attachments") {
		return facades.Schema().Create("attachments", func(table schema.Blueprint) {
			table.ID()
			table.BigInteger("episode_id").Default(0) // 所属课程 ID
			table.String("name")                      // 附件名称
			table.String("file_path")                 // 附件路径
			table.TimestampsTz()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20260318162447CreateAttachmentsTable) Down() error {
	return facades.Schema().DropIfExists("attachments")
}
