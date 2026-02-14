package middleware

import (
	"homecourse/app/http/response"
	"homecourse/app/models"

	"github.com/goravel/framework/contracts/http"
)

func Admin() http.Middleware {
	return func(ctx http.Context) {
		cuser := ctx.Value(models.Cuser).(models.User)

		if cuser.Role != models.RoleAdmin {
			_ = response.Unauthorized(ctx, "E5").Abort()
			return
		}

		ctx.Request().Next()
	}
}
