package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"homecourse/app/facades"
)

type M20260209201058CreateEpisodesTable struct{}

// Signature The unique signature for the migration.
func (r *M20260209201058CreateEpisodesTable) Signature() string {
	return "20260209201058_create_episodes_table"
}

// Up Run the migrations.
func (r *M20260209201058CreateEpisodesTable) Up() error {
	if !facades.Schema().HasTable("episodes") {
		return facades.Schema().Create("episodes", func(table schema.Blueprint) {
			table.ID()                                   // 自增 ID
			table.BigInteger("group_id").Default(0)      // 所属分组 ID
			table.String("title")                        // 剧集标题（如：01. 环境搭建）
			table.String("file_path")                    // 视频文件的路径
			table.Integer("sort").Default(0)             // 按课程分组排序
			table.Integer("duration").Default(0)         // 时长（秒）（后续可用 ffmpeg 自动获取）
			table.Integer("watched_duration").Default(0) // 已观看时长（秒）
			table.TimestampsTz()

			table.Foreign("group_id").References("id").On("groups").CascadeOnDelete()

			table.Index("group_id")
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20260209201058CreateEpisodesTable) Down() error {
	return facades.Schema().DropIfExists("episodes")
}
