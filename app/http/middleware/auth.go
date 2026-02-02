package middleware

import (
	"homecourse/app/http/response"
	"homecourse/app/models"
	"strconv"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

func Auth() http.Middleware {
	return func(ctx http.Context) {
		token := ctx.Request().Header("Authorization")

		if token == "" || !strings.HasPrefix(token, "Bearer ") {
			_ = response.Unauthorized(ctx, "E1").Abort()
			return
		}
		token = strings.TrimPrefix(token, "Bearer ")

		if _, err := facades.Auth(ctx).Parse(token); err != nil {
			_ = response.Unauthorized(ctx, "E2").Abort()
			return
		}

		// 获取登录用户ID
		uid, err := facades.Auth(ctx).ID()
		if err != nil {
			_ = response.Unauthorized(ctx, "E3").Abort()
			return
		}
		ok, err := strconv.ParseUint(uid, 10, 0)
		if err != nil {
			_ = response.Unauthorized(ctx, "E4").Abort()
			return
		}

		// 判断此id是否存在
		if flag, err := facades.Orm().Query().Model(&models.User{}).Where("id", uid).Exists(); err != nil {
			_ = response.Unauthorized(ctx, "E5").Abort()
			return
		} else if !flag {
			_ = response.Unauthorized(ctx, "E6").Abort()
			return
		}

		// 存放UserID
		ctx.WithValue(models.UserID, uint(ok))

		ctx.Request().Next()
	}
}
