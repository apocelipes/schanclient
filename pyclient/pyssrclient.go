package pyclient

import (
	"os/exec"
	"os"

	"schanclient/config"
	"schanclient/ssr"
)

// Python ssr客户端
type PySSRClient struct {
	bin    string
	binArg string
	config string
}

func init() {
	ssr.SetLuancherMaker("python", ssr.LauncherMaker(NewPySSRClient))
}

func NewPySSRClient(c *config.UserConfig) ssr.SSRLauncher {
	p := new(PySSRClient)
	bin, err := c.SSRBin.AbsPath()
	if err != nil {
		return nil
	}
	p.bin = bin
	p.config, err = c.SSRConfigPath.AbsPath()
	if err != nil {
		return nil
	}

	p.binArg = "-c" + p.config

	return p
}

func (p *PySSRClient) Start() error {
	// env tells sudo to find commands in $PATH
	cmd := exec.Command("sudo", "env PATH=$PATH", "python", p.bin, p.binArg, "-d", "start")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (p *PySSRClient) Restart() error {
	cmd := exec.Command("sudo", "env PATH=$PATH", "python", p.bin, p.binArg, "-d", "restart")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (p *PySSRClient) Stop() error {
	cmd := exec.Command("sudo", "env PATH=$PATH", "python", p.bin, p.binArg, "-d", "stop")
	cmd.Stderr = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
