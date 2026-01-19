package errcode

const (
	Success     = 0
	ServerError = 10001
	ParamError  = 10002
	NotFound    = 10003

	// User (20000 - 29999)
	UserExists   = 20001
	UserNotFound = 20002
	AuthFailed   = 20003
	TokenInvalid = 20004

	// Article (30000 - 39999)
	ArticleNotFound = 30001
)

// TODO: International sufficiency
// code msg
var msgFlags = map[int]string{
	Success:     "ok",
	ServerError: "系统内部错误，请稍后再试",
	ParamError:  "请求参数错误",
	NotFound:    "资源不存在",

	UserExists:   "用户已存在",
	UserNotFound: "用户不存在",

	AuthFailed:   "用户名或密码错误",
	TokenInvalid: "登录已过期，请重新登录",

	ArticleNotFound: "文章不存在",
}

func GetMsg(code int) string {
	msg, ok := msgFlags[code]
	if ok {
		return msg
	}
	return msgFlags[ServerError]
}
