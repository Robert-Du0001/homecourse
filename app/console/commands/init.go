package commands

import (
	"homecourse/app/facades"
	"homecourse/app/models"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
)

type Init struct {
}

// Signature The name and signature of the console command.
func (r *Init) Signature() string {
	return "app:init"
}

// Description The console command description.
func (r *Init) Description() string {
	return "初始化项目数据"
}

// Extend The console command extend.
func (r *Init) Extend() command.Extend {
	return command.Extend{Category: "app"}
}

// Handle Execute the console command.
func (r *Init) Handle(ctx console.Context) error {
	// 处理课程分类数据
	if exists, err := facades.Orm().Query().
		Model(&models.Category{}).
		Exists(); err != nil {
		facades.Log().Error("初始化课程分类数据失败[E1]", err)
		ctx.Error("初始化课程分类数据失败[E2]: " + err.Error())
	} else if !exists {
		// 若没有课程分类数据，需要添加个默认分类
		category := models.Category{
			Name:      "默认分类",
			IsDefault: true,
		}
		if err := facades.Orm().Query().Create(&category); err != nil {
			facades.Log().Error("初始化课程分类数据失败[E2]", err)
			ctx.Error("初始化课程分类数据失败[E2]: " + err.Error())
		}
	}

	ctx.Success("初始化成功")
	return nil
}
