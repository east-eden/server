package auto

import "errors"

// ErrGlobalConfigInvalid 全局配置缺失
var ErrGlobalConfigInvalid = errors.New("invalid global config")

// GetGlobalConfig 返回全局配置（GlobalConfig.csv 第一行）
func GetGlobalConfig() (*GlobalConfigEntry, bool) {
	return GetGlobalConfigEntry(1)
}
