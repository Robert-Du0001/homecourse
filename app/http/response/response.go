package response

import "github.com/goravel/framework/contracts/http"

// 请求成功 - 200
func Ok(ctx http.Context, msg string, data any) http.Response {
	return ctx.Response().Json(
		http.StatusOK,
		http.Json{
			"msg":  msg,
			"data": data,
		},
	)
}

// 参数错误 - 400
func BadRequest(ctx http.Context, msg string, data any) http.Response {
	return ctx.Response().Json(
		http.StatusBadRequest,
		http.Json{
			"msg":  msg,
			"data": data,
		},
	)
}

// 身份验证失败 - 401
func Unauthorized(ctx http.Context, code string) http.AbortableResponse {
	return ctx.Response().Json(
		http.StatusUnauthorized,
		http.Json{
			"msg":  "身份验证失败[" + code + "]",
			"data": nil,
		},
	)
}

// 服务端错误 - 500
func InternalServerError(ctx http.Context, code string, err error) http.Response {
	return ctx.Response().Json(
		http.StatusInternalServerError,
		http.Json{
			"msg":  err.Error() + "[" + code + "]",
			"data": nil,
		},
	)
}
