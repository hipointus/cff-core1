package mobile

import (
	"os"
	"path/filepath"

	"github.com/metacubex/mihomo/constant"
	"github.com/metacubex/mihomo/hub/executor"
	"github.com/metacubex/mihomo/hub/route"
	"github.com/metacubex/mihomo/listener"
	"github.com/metacubex/mihomo/log"
	"github.com/metacubex/mihomo/tunnel"
	"github.com/oschwald/geoip2-golang"
	"go.uber.org/automaxprocs/maxprocs"
)

// status service status
var status = false

func init() {
	_, _ = maxprocs.Set(maxprocs.Logger(func(string, ...any) {}))
}

func SetHomeDir(homeDir string) bool {
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

func SetConfig(configFile string) bool {
	if configFile == "" {
		return false
	}
	if !filepath.IsAbs(configFile) {
		configFile = filepath.Join(constant.Path.HomeDir(), configFile)
	}
	constant.SetConfig(configFile)
	return true
}

func VerifyMMDB(path string) bool {
	instance, err := geoip2.Open(path)
	if err == nil {
		_ = instance.Close()
	}
	return err == nil
}

func StartRust(addr string) string {
	go route.Start(addr, "")
	oldAddr := route.GetAddr()
	if oldAddr == "" {
		return addr
	}
	return oldAddr
}

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

func OperateTun(enable bool, fileDescriptor, mtu int32) {
	tun := listener.LastTunConf
	tun.Enable = enable
	tun.MTU = uint32(mtu)
	tun.FileDescriptor = int(fileDescriptor)
	tun.TunIf = true
	listener.ReCreateTun(tun, tunnel.TCPIn(), tunnel.UDPIn())
}
