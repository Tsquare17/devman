package environment

import (
	"os/exec"
)

var nginxName = "nginx"
var apacheName = "apache2"
var caddyName = "caddy"

func IsRunningNginx() bool {
	return isRunningProcess(nginxName)
}

func IsRunningApache() bool {
	return isRunningProcess(apacheName)
}

func IsRunningCaddy() bool {
	return isRunningProcess(caddyName)
}

func WebServerProcessName() string {
	if IsRunningApache() {
		return apacheName
	} else if IsRunningNginx() {
		return nginxName
	} else if IsRunningCaddy() {
		return caddyName
	} else {
		return ""
	}
}

func isRunningProcess(name string) bool {
	cmd := exec.Command("pidof", name)
	err := cmd.Run()
	if err != nil {
		return false
	}

	return true
}
