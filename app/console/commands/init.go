package commands

import (
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

	ctx.Success("初始化成功")
	return nil
}
