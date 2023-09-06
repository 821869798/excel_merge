//go:build !windows

package reg

import "errors"

func ReadRegistry(wk WindowsRegistryKey, path, key string) (string, error) {
	return "", errors.New("not support in current system")
}

func WriteRegistry(wk WindowsRegistryKey, path, key, value string) error {
	return errors.New("not support in current system")
}
