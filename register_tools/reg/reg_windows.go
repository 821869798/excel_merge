//go:build windows

package reg

import (
	"golang.org/x/sys/windows/registry"
)

func ReadRegistry(k registry.Key, path, key string) (string, error) {
	// 打开注册表项
	k, err := registry.OpenKey(k, path, registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}
	defer k.Close()

	// 获取注册表键
	value, _, err := k.GetStringValue(key)
	if err != nil {
		return "", err
	}
	return value, nil
}

func WriteRegistry(k registry.Key, path, key, value string) error {
	k, err := registry.OpenKey(k, path, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()

	err = k.SetStringValue(key, value)
	return err
}
