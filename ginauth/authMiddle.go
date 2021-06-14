package ginauth

import (
	"encoding/base64"
	"encoding/json"
	"github.com/devloperPlatform/go-base-utils/commonvos"
	"github.com/devloperPlatform/go-gin-base/ginbase"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	HeaderUserFlag    = "u_token"
	HeaderSysFlag     = "bk_sys"
	HeaderDevUserFlag = "bk_dev_user"
)

type UserAuthFun func(user *commonvos.InsideUserInfo, ctx *gin.Context) interface{}

type AuthMiddle struct {
	*gin.Engine
	dev      bool
	tokenUrl string
	//querySysFlag         func(context *gin.Context) string
	//sysFlag              string
	needHandleFn         func(ctx *gin.Context) bool
	convertUserInterface ConvertUserInterface
}

// ConvertUserInterface 转换用户接口方法
type ConvertUserInterface func(data []byte) (interface{}, error)

func (a *AuthMiddle) StartDev() *AuthMiddle {
	a.dev = true
	return a
}

func (a *AuthMiddle) SettingTokenUrl(tokenUrl string) *AuthMiddle {
	a.tokenUrl = tokenUrl
	return a
}

func (a *AuthMiddle) SettingConvertUserFun(convertUserFun ConvertUserInterface) *AuthMiddle {
	a.convertUserInterface = convertUserFun
	return a
}

//func (a *AuthMiddle) SettingSysFlagFn(fn func(context *gin.Context) string) *AuthMiddle {
//	a.querySysFlag = fn
//	return a
//}

//func (a *AuthMiddle) SettingSysFlag(sys string) *AuthMiddle {
//	a.sysFlag = sys
//	return a
//}

func (a *AuthMiddle) SettingNeedHandleFn(fn func(ctx *gin.Context) bool) *AuthMiddle {
	a.needHandleFn = fn
	return a
}

func (a *AuthMiddle) AuthControlMiddle() gin.HandlerFunc {
	return func(context *gin.Context) {
		httpResultUtil := ginbase.HttpResultUtil()
		isServerUrl := false
		if a.needHandleFn != nil {
			isServerUrl = a.needHandleFn(context)
		}

		header := context.GetHeader(HeaderUserFlag)
		//sysHeader := context.GetHeader(HeaderSysFlag)
		if header == "" {
			if a.dev || isServerUrl {
				userInterface, err := a.convertUserInterface(nil)
				if err == nil {
					context.Set("nowUser", userInterface)
				}
				context.Next()
				return
			}

			context.JSON(401, httpResultUtil.NewErrorHttpResultWithCodeAndMsg("401", "用户认证失败"))
			context.Abort()
			return
		}

		//sysHeader := a.sysFlag
		//if a.querySysFlag != nil {
		//	sysHeader = a.querySysFlag(context)
		//}

		//if sysHeader != "" {
		//	sysHeader = base64.StdEncoding.EncodeToString([]byte(sysHeader))
		response, err := http.PostForm(a.tokenUrl, map[string][]string{
			"token": {header},
		})

		var data []byte
		if err != nil {
			if !a.dev {
				context.JSON(401, httpResultUtil.NewErrorHttpResultWithCodeAndMsg("401", err.Error()))
				context.Abort()
				return
			}
			devUserInfo := context.GetHeader(HeaderDevUserFlag)
			if devUserInfo == "" {
				context.Next()
				return
			}
			decodeString, err := base64.StdEncoding.DecodeString(devUserInfo)
			if err != nil {
				context.Next()
				return
			}
			data = decodeString
		} else {
			defer response.Body.Close()
			all, err := ioutil.ReadAll(response.Body)
			if err != nil {
				context.JSON(500, httpResultUtil.NewErrorHttpResultWithCodeAndMsg("401", "读取认证"))
				context.Abort()
				return
			}

			var res *ginbase.HttpResult
			err = json.Unmarshal(all, &res)
			if err != nil {
				panic("解析数据失败")
			}
			if res.Error {
				context.JSON(401, httpResultUtil.NewErrorHttpResultWithMsg(res.Msg))
				context.Abort()
				return
			}
			data = all
		}

		if a.convertUserInterface != nil {
			userInterface, err := a.convertUserInterface(data)
			if err != nil {
				context.JSON(500, gin.H{
					"error": true,
					"code":  "500",
					"msg":   err.Error(),
				})
				context.Abort()
				return
			}
			context.Set("nowUser", userInterface)
			context.Next()
			return
		}

		context.Abort()
		panic("非法用户访问!")
		//}
	}
}

func (a *AuthMiddle) WrapperRequestWithUserInfoResponseInterface(fn UserAuthFun) func(ctx *gin.Context) {
	if fn == nil {
		panic("处理函数不能为空")
	}
	return ginbase.HttpResultUtil().WrapperResponseHandle(func(ctx *gin.Context) interface{} {
		nowUser, exists := ctx.Get("nowUser")
		if !exists {
			panic("获取当前登录用户失败!")
		}
		if info, ok := nowUser.(commonvos.InsideUserInfo); ok {
			return fn(&info, ctx)

		}

		if info, ok := nowUser.(*commonvos.InsideUserInfo); ok {
			return fn(info, ctx)

		}
		panic("解析用户信息失败!")
	})
}

func (a *AuthMiddle) GET(relativePath string, fun UserAuthFun) *AuthMiddle {
	a.Engine.GET(relativePath, a.WrapperRequestWithUserInfoResponseInterface(fun))
	return a
}

func (a *AuthMiddle) POST(relativePath string, fun UserAuthFun) *AuthMiddle {
	a.Engine.POST(relativePath, a.WrapperRequestWithUserInfoResponseInterface(fun))
	return a
}

func (a *AuthMiddle) DELETE(relativePath string, fun UserAuthFun) *AuthMiddle {
	a.Engine.DELETE(relativePath, a.WrapperRequestWithUserInfoResponseInterface(fun))
	return a
}

func (a *AuthMiddle) PATCH(relativePath string, fun UserAuthFun) *AuthMiddle {
	a.Engine.PATCH(relativePath, a.WrapperRequestWithUserInfoResponseInterface(fun))
	return a
}

func (a *AuthMiddle) PUT(relativePath string, fun UserAuthFun) *AuthMiddle {
	a.Engine.PUT(relativePath, a.WrapperRequestWithUserInfoResponseInterface(fun))
	return a
}

func (a *AuthMiddle) OPTIONS(relativePath string, fun UserAuthFun) *AuthMiddle {
	a.Engine.OPTIONS(relativePath, a.WrapperRequestWithUserInfoResponseInterface(fun))
	return a
}

func (a *AuthMiddle) HEAD(relativePath string, fun UserAuthFun) *AuthMiddle {
	a.Engine.HEAD(relativePath, a.WrapperRequestWithUserInfoResponseInterface(fun))
	return a
}

func New(engine *gin.Engine) *AuthMiddle {
	return &AuthMiddle{
		engine,
		false,
		os.Getenv("AUTH_TOKEN_URL"),
		nil,
		func(data []byte) (interface{}, error) {
			var user *commonvos.InsideUserInfo
			err := json.Unmarshal(data, &user)
			return user, err
		},
	}
}

func EnabledAuth(engine *gin.Engine) *AuthMiddle {
	middle := New(engine)
	middle.Use(middle.AuthControlMiddle())
	return middle
}
