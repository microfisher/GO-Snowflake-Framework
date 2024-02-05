package example

import (
	"context"
	"database/sql"
	"fmt"
	"snowflake/comm"
	"snowflake/config"
	"snowflake/core"
	"time"

	"github.com/redis/go-redis/v9"
)

// ExampleServer .
type ExampleServer struct {
	ctx                context.Context    // 上下文信息
	ctxCancel          context.CancelFunc // 全局控制
	mysql              *sql.DB            // 数据库
	redis              *redis.Client      // 缓存
	wssEndpoints       []string           // 节点列表
	httpsEndpoints     []string           // 节点列表
	appIntervalSeconds int                // 运行环境报告时间
	webListen          string             // 网站监听
	webHeaders         string             // 网站头部
	webDirectory       map[string]string  // 网站目录
	webJwtSecret       string             // 网站密钥
	webSignature       string             // 网站签名
	webAllowOrigins    string             // 网站起源
	module             string             // 模块名称
	timeout            int                // 超时时间
	debug              bool               // 是否测试
}

// NewExampleServer .
func NewExampleServer() comm.Runable {
	ctx, cancel := context.WithCancel(context.Background())
	data := &ExampleServer{
		ctx:       ctx,
		ctxCancel: cancel,
		module:    comm.GetPackageName(NewExampleServer),
	}
	return data
}

// Start .
func (s *ExampleServer) Start() {
	if !s.CheckData() {
		return
	}
	if !s.CheckConfig() {
		return
	}

	go s.goTicker()
	go s.goWebServer()

}

// Stop .
func (s *ExampleServer) Stop() {
	s.ctxCancel()
	s.Info("server has stoped.")
}

// Done .
func (s *ExampleServer) Done() <-chan struct{} {
	return s.ctx.Done()
}

// ControlCmd 外部控制api传入命令，解析之后推入事件队列
func (s *ExampleServer) ControlCmd(body string) (string, error) {
	return "unsupport", nil
}

// 信息日志
func (s *ExampleServer) Info(format string, args ...any) {
	comm.Info(s.module, format, args...)
}

// 错误日志
func (s *ExampleServer) Error(format string, args ...any) {
	comm.Error(s.module, format, args...)
}

// 初始化数据库
func (s *ExampleServer) CheckData() bool {
	mysql := core.OpenMysql(s.ctx, s.module)
	if mysql == nil {
		return false
	}
	s.mysql = mysql

	redis := core.OpenRedis(s.ctx, s.module)
	if redis == nil {
		return false
	}
	s.redis = redis

	return true
}

// 初始化配置
func (s *ExampleServer) CheckConfig() bool {

	s.wssEndpoints = config.GetViper().GetStringSlice("wss_endpoints")
	if len(s.wssEndpoints) <= 0 {
		s.Error("failed to get wss endpoints.")
		return false
	}

	s.httpsEndpoints = config.GetViper().GetStringSlice("https_endpoints")
	if len(s.httpsEndpoints) <= 0 {
		s.Error("failed to get http endpoints.")
		return false
	}

	s.webListen = config.GetViper().GetString(fmt.Sprintf("%s.web_listen", s.module))
	if len(s.webListen) <= 0 {
		s.Error("failed to get web listen.")
		return false
	}

	s.webHeaders = config.GetViper().GetString(fmt.Sprintf("%s.web_headers", s.module))
	if len(s.webHeaders) <= 0 {
		s.Error("failed to get web headers.")
		return false
	}

	s.webDirectory = config.GetViper().GetStringMapString(fmt.Sprintf("%s.web_directory", s.module))
	if len(s.webDirectory) <= 0 {
		s.Error("failed to get web directory.")
		return false
	}

	s.webJwtSecret = config.GetViper().GetString(fmt.Sprintf("%s.web_jwt_secret", s.module))
	if len(s.webJwtSecret) <= 0 {
		s.Error("failed to get web jwt secret.")
		return false
	}

	s.webSignature = config.GetViper().GetString(fmt.Sprintf("%s.web_signature", s.module))
	if len(s.webSignature) <= 0 {
		s.Error("failed to get web signature.")
		return false
	}

	s.webAllowOrigins = config.GetViper().GetString(fmt.Sprintf("%s.web_allow_origins", s.module))
	if len(s.webAllowOrigins) <= 0 {
		s.Error("failed to get web allow origins.")
		return false
	}

	s.timeout = config.GetViper().GetInt(fmt.Sprintf("%s.timeout_seconds", s.module))
	if s.timeout <= 0 {
		s.timeout = 30
	}

	s.appIntervalSeconds = config.GetViper().GetInt((fmt.Sprintf("%s.app_interval_seconds", s.module)))
	if s.appIntervalSeconds <= 0 {
		s.appIntervalSeconds = 60
	}

	s.debug = config.GetViper().GetBool(s.module + ".debug")

	return true
}

// 定时打印日志
func (s *ExampleServer) goTicker() {
	defer s.Recover()
	timer := time.NewTicker(time.Second * time.Duration(s.appIntervalSeconds))
	defer timer.Stop()
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-timer.C:
			memory, goroutine := comm.GetSystemInfo()
			s.Info("app runtime status: %dmb -> %d", memory, goroutine)
		}
	}
}

// 协程防崩溃
func (s *ExampleServer) Recover(args ...any) {
	if r := recover(); r != nil {
		s.Error("failed to execute goroutine -> args:%v r:%v caller_stack:%s", args, r, comm.GetCallerStackLog())
	}
}
