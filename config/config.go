package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

const (
	// 默认配置文件路径
	configPath = ".local/share/schanclient.json"
)

var (
	// 无法查找$HOME
	ErrHOME = errors.New("can't find $HOME in your environments.")
	// 路径无法解析为绝对路径
	ErrNotAbs = errors.New("path is not an abs path")
)

// 用户配置
type UserConfig struct {
	UserName string `json:"user_name"`
	Passwd   string `json:"user_password"`
	// ssr config
	SSRConfigPath JSONPath `json:"ssr_config_path"`
	// ssr client bin path
	SSRBin  JSONPath `json:"ssr_bin"`
	LogFile JSONPath `json:"log_file"`
}

// ConfigPath 返回`～`被替换为$HOME的config path
func ConfigPath() (string, error) {
	home, exist := os.LookupEnv("HOME")
	if !exist {
		return "", ErrHOME
	}

	return home + string(os.PathSeparator) + configPath, nil
}

func (u *UserConfig) StoreConfig() error {
	storePath, err := ConfigPath()
	if err != nil {
		return err
	}

	f, err := os.OpenFile(storePath, os.O_WRONLY|os.O_CREATE, 0664)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := json.MarshalIndent(u, "", "\t")
	if err != nil {
		return err
	}

	if _, err = f.Write(data); err != nil {
		return err
	}

	return nil
}

func (u *UserConfig) LoadConfig() error {
	loadPath, err := ConfigPath()
	if err != nil {
		return err
	}

	f, err := os.Open(loadPath)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(data, u); err != nil {
		return err
	}

	return nil
}
