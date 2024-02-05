package comm

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"snowflake/config"
	"snowflake/log"
)

// ControlService 支持控制命令的服务
type ControlService interface {
	ControlCmd(body string) (string, error)
}

// GoControlListen 启动命令调用接口
func GoControlListen(svr ControlService, listen string) {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("event:%s msg:control_server_exception %v caller_stack:%s",
				EventWarn, r, GetCallerStackLog())
		}
	}()
	if len(listen) == 0 {
		log.Infof("control port not enabled")
		return
	}
	log.Infof("control port listen at:%s", listen)
	mux := http.NewServeMux()
	mux = DefaultControlHandler(svr, mux)
	err := http.ListenAndServe(listen, mux)
	if err != nil {
		log.Errorf("event:%s msg:controller_at %s err %s", EventWarn, listen, err.Error())
		return
	}
}

// DefaultControlHandler .
func DefaultControlHandler(s ControlService, mux *http.ServeMux) *http.ServeMux {
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})
	mux.HandleFunc("/reloadconfig", func(w http.ResponseWriter, r *http.Request) {
		err := config.ReloadConfig("json")
		if err != nil {
			w.Write([]byte(fmt.Sprintf("reload config err:%s", err.Error())))
			return
		}
		w.Write([]byte(fmt.Sprintf("reload config success")))
	})
	mux.HandleFunc("/cmd", func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			log.Errorf("/cmd read body err:%s", err.Error())
			w.Write([]byte(err.Error()))
			return
		}
		result, err := s.ControlCmd(string(data))
		w.Write([]byte(result))
	})
	return mux
}

// PostControlCmd 发送一个控制命令到端口
func PostControlCmd(controlListen, cmdBody string) ([]byte, error) {
	url := fmt.Sprintf("http://%s/cmd", controlListen)
	bodyReader := bytes.NewReader([]byte(cmdBody))
	resp, err := http.Post(url, "text/plain", bodyReader)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	return respBody, nil
}

// PostReloadConfigCmd .
func PostReloadConfigCmd(controlListen string) ([]byte, error) {
	url := fmt.Sprintf("http://%s/reloadconfig", controlListen)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	return respBody, nil
}
