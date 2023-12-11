package main

import "C"
import (
	"os"
	"path/filepath"

	"github.com/metacubex/mihomo/constant"
	"github.com/metacubex/mihomo/hub/executor"
	"github.com/metacubex/mihomo/hub/route"
	"github.com/metacubex/mihomo/log"
	"github.com/oschwald/geoip2-golang"
	"go.uber.org/automaxprocs/maxprocs"
)

// status service status
var status = false

func init() {
	_, _ = maxprocs.Set(maxprocs.Logger(func(string, ...any) {}))
}

//export SetHomeDir
func SetHomeDir(homeStr *C.char) bool {
	homeDir := C.GoString(homeStr)
	info, err := os.Stat(homeDir)
	if err != nil {
		log.Errorln("[Clash Lib] SetHomeDir: %s : %+v", homeDir, err)
		return false
	}
	if !info.IsDir() {
		log.Errorln("[Clash Lib] SetHomeDir: Path is not directory %s", homeDir)
		return false
	}
	constant.SetHomeDir(homeDir)
	return true
}

//export SetConfig
func SetConfig(configStr *C.char) bool {
	configFile := C.GoString(configStr)
	if configFile == "" {
		return false
	}
	if !filepath.IsAbs(configFile) {
		configFile = filepath.Join(constant.Path.HomeDir(), configFile)
	}
	constant.SetConfig(configFile)
	return true
}

//export VerifyMMDB
func VerifyMMDB(path *C.char) bool {
	instance, err := geoip2.Open(C.GoString(path))
	if err == nil {
		_ = instance.Close()
	}
	return err == nil
}

//export StartRust
func StartRust(addr *C.char) *C.char {
	go route.Start(C.GoString(addr), "")
	oldAddr := route.GetAddr()
	if oldAddr == "" {
		return addr
	}
	return C.CString(oldAddr)
}

//export StartService
func StartService() bool {
	if status {
		return status
	}

	if constant.Path.Config() == "config.yaml" {
		configFile := filepath.Join(constant.Path.HomeDir(), constant.Path.Config())
		constant.SetConfig(configFile)
	}

	cfg, err := executor.Parse()
	if err != nil {
		log.Errorln("[Clash Lib] StartService: Parse config error: %+v", err)
		return status
	}
	executor.ApplyConfig(cfg, true)

	status = true
	return status
}

func main() {}
