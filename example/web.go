package example

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/golang-jwt/jwt/v4"
)

func (s *ExampleServer) goWebServer() {

	// 限制上传
	app := fiber.New(fiber.Config{
		BodyLimit: 2 * 1024 * 1024, // 2MB
	})

	// 跨域访问
	app.Use(cors.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: s.webAllowOrigins,
		AllowHeaders: s.webHeaders,
	}))

	// 限制频率
	app.Use(limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == "127.0.0.1"
		},
		Max:        60,
		Expiration: 60 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Get("x-forwarded-for")
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.JSON(ResultObject{Code: 0, Data: "failed to process your request: reached the max limit"})
		},
	}))

	// 静态目录
	for key, value := range s.webDirectory {
		app.Static(key, value)
	}

	// 系统接口（无需验证）
	app.Post("/api/user/signin", s.signIn)

	// 系统鉴权
	// app.Use(jwtware.New(jwtware.Config{
	// 	SigningKey: []byte(s.jwtSecret),
	// }))

	// 系统接口（需要验证）
	// app.Post("/api/painter/draw", s.drawImage)
	// app.Post("/api/painter/upload", s.uploadImage)
	// app.Get("/api/painter/prompts/:id?/:page?/:size?", s.getPrompts)
	// app.Get("/api/painter/task/latest", s.getTaskInfo)

	// 启动服务
	s.Info("listening on: %s", s.webListen)
	err := app.Listen(s.webListen)
	if err != nil {
		s.Error("failed to listening on: %s -> %s", s.webListen, err.Error())
		return
	}

}

// 用户登陆
func (s *ExampleServer) signIn(c *fiber.Ctx) error {

	data := ResultObject{Code: 0}
	username := c.FormValue("username")
	password := c.FormValue("password")

	if !VerifySignature(username, s.webSignature, password) {
		data.Code = 1
		data.Data = "Failed to verify your username or password"
		return c.JSON(data)
	}

	claims := jwt.MapClaims{
		"id":     1,
		"name":   username,
		"level":  1,
		"expire": time.Now().Add(time.Hour * 168).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signature, err := token.SignedString([]byte(s.webJwtSecret))
	if err != nil {
		data.Code = 2
		data.Data = err.Error()
		return c.JSON(data)
	}

	data.Data = signature
	return c.JSON(data)
}

// 获取请求数据
func GetRequestData(ctx *fiber.Ctx) string {

	var request = ""
	request += fmt.Sprintf("Request method is %s\n", ctx.Method())
	request += fmt.Sprintf("RequestURI is %s\n", ctx.OriginalURL())
	request += fmt.Sprintf("Requested path is %s\n", ctx.Path())
	request += fmt.Sprintf("Host is %s\n", ctx.Hostname())
	request += fmt.Sprintf("Query string is %s\n", string(ctx.Context().QueryArgs().QueryString()))
	request += fmt.Sprintf("User-Agent is %s\n", string(ctx.Context().UserAgent()))
	request += fmt.Sprintf("Request call time is %s\n", ctx.Context().Time().String())
	request += fmt.Sprintf("Serial request number for the current connection is %d\n", ctx.Context().ConnRequestNum())
	request += fmt.Sprintf("Your ip is %q\n\n", ctx.Context().RemoteIP().String())
	request += fmt.Sprintf("Body is:\n%s\n", string(ctx.Body()))

	return request
}
