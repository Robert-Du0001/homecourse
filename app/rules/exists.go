package rules

import (
	"context"
	"homecourse/app/facades"

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
	optionsLen := len(options)

	if optionsLen < 1 {
		return false
	} else if optionsLen == 1 {
		options = append(options, "id")
	}

	if exists, err := facades.Orm().Query().Table(options[0].(string)).
		Where(options[1].(string), val.(uint)).
		Exists(); err != nil {
		facades.Log().Error(err)
		return false
	} else if !exists {
		return false
	}

	return true
}

// Message Get the validation error message.
func (receiver *Exists) Message(ctx context.Context) string {
	// exists:categories,category_id
	return ":attribute 字段不存在值"
}
