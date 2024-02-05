package debug

import (
	"net/http"

	"snowflake/comm"
	"snowflake/log"
)

// GoPprofListen 监控服务
func GoPprofListen(listen string) {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("event:%s msg:pprof_server_exception %v caller_stack:%s",
				comm.EventWarn, r, comm.GetCallerStackLog())
		}
	}()
	if len(listen) == 0 {
		log.Infof("pprof port not enabled")
		return
	}
	log.Infof("pprof port listen at:%s", listen)
	log.Error("pprof error:", http.ListenAndServe(listen, nil))
}
