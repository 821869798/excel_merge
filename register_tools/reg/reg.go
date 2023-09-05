//go:build !windows

package reg

import "errors"

func ReadRegistry(k registry.Key, path, key string) (string, error) {
	return "", errors.New("not support in current system")
}

func WriteRegistry(k registry.Key, path, key, value string) error {
	return errors.New("not support in current system")
}
