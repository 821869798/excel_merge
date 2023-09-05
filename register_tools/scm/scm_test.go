package scm

import (
	"os"
	"testing"
)

func TestAppData(t *testing.T) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		t.Errorf("无法获取用户的缓存目录:%v", err)
		return
	}
	t.Logf("path:%v", cacheDir)

}
