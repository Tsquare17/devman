package user

import (
	"os"
	"os/user"
	"syscall"
)

func IsSuper() bool {
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}

	if currentUser.Uid == "0" {
		return true
	}

	return false
}

func GetName() string {
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}

	if currentUser.Uid != "0" {
		return currentUser.Username
	}

	return os.Getenv("SUDO_USER")
}

func GetRegularUserName() string {
	return os.Getenv("SUDO_USER")
}

func GetUidGid(file string) (int, int) {
	info, _ := os.Stat(file)

	var uid int
	var gid int
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		uid = int(stat.Uid)
		gid = int(stat.Gid)
	} else {
		uid = os.Getuid()
		gid = os.Getgid()
	}

	return uid, gid
}
