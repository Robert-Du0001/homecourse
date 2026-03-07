package rules

import (
	"context"
	"fmt"

	"github.com/goravel/framework/contracts/validation"
)

type Exists struct {
}

// Signature The name of the rule.
func (receiver *Exists) Signature() string {
	return "exists"
}

// Passes Determine if the validation rule passes.
func (receiver *Exists) Passes(ctx context.Context, data validation.Data, val any, options ...any) bool {
	fmt.Println(data, val, options)

	return false
}

// Message Get the validation error message.
func (receiver *Exists) Message(ctx context.Context) string {
	// exists:categories,category_id
	return "不存在 :attribute 值"
}
