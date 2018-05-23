package ssr

import (
	"schanclient/config"
)

type SSRLauncher interface {
	Start() error
	Restart() error
	Stop() error
}

type LauncherMaker func(*config.UserConfig) SSRLauncher

var launchers = make(map[string]LauncherMaker)

func SetLuancherMaker(name string, l LauncherMaker) {
	if name == "" || l == nil {
		panic("SetLauncher error: wrong name or LuancherMaker.")
	}
	
	launchers[name] = l
}

func NewLauncher(name string, conf *config.UserConfig) SSRLauncher {
	l, ok := launchers[name]
	if !ok {
		return nil
	}
	
	return l(conf)
}
