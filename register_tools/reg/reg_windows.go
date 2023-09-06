//go:build windows

package reg

import (
	"golang.org/x/sys/windows/registry"
)

func GetRegistryKey(k WindowsRegistryKey) registry.Key {
	switch k {
	case CLASSES_ROOT:
		return registry.CLASSES_ROOT
	case CURRENT_USER:
		return registry.CURRENT_USER
	case LOCAL_MACHINE:
		return registry.LOCAL_MACHINE
	case USERS:
		return registry.USERS
	case CURRENT_CONFIG:
		return registry.CURRENT_CONFIG
	case PERFORMANCE_DATA:
		return registry.PERFORMANCE_DATA
	default:
		return registry.CLASSES_ROOT
	}
}

func ReadRegistry(wk WindowsRegistryKey, path, key string) (string, error) {
	k := GetRegistryKey(wk)
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

func WriteRegistry(wk WindowsRegistryKey, path, key, value string) error {
	k := GetRegistryKey(wk)
	k, err := registry.OpenKey(k, path, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()

	err = k.SetStringValue(key, value)
	return err
}
