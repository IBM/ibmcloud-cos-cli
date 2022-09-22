package aspera

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSDKDir(t *testing.T) {
	os.Setenv("ASPERA_SDK_PATH", "/path/to/sdk")
	assert.Equal(t, "/path/to/sdk", SDKDir())
	os.Unsetenv("ASPERA_SDK_PATH")

	if home, err := os.UserHomeDir(); err == nil {
		assert.Equal(t, filepath.Join(home, ".aspera_sdk"), SDKDir())
	} else {
		assert.Equal(t, ".aspera_sdk", SDKDir())
	}
}

func TestTransferdBinPath(t *testing.T) {
	os.Setenv("ASPERA_SDK_PATH", "/path/to/sdk")
	if runtime.GOOS == "windows" {
		assert.Equal(t, "/path/to/sdk/bin/asperatransferd.exe", TransferdBinPath())
	} else {
		assert.Equal(t, "/path/to/sdk/bin/asperatransferd", TransferdBinPath())
	}
}
func TestGetSDKAttributes(t *testing.T) {
	supportedPlatforms := []struct {
		os       string
		arch     string
		platform string
		ext      string
	}{
		{"darwin", "amd64", "osx-amd64", "tar.gz"},   //special os
		{"windows", "amd64", "windows-amd64", "zip"}, //special ext
		{"linux", "amd64", "linux-amd64", "tar.gz"},  //common case
	}

	for _, pair := range supportedPlatforms {
		platform, ext, _ := getSDKAttributes(pair.os, pair.arch)
		assert.Equal(t, platform, pair.platform)
		assert.Equal(t, ext, pair.ext)
	}
}

func TestGetSDKAttributesError(t *testing.T) {
	notSupportedPlatforms := []struct {
		os   string
		arch string
	}{
		{"darwin", "arm64"},  // not supported os
		{"freebsd", "amd64"}, // not supported arch
		{"orbis", "amd64"},   // either is not supported
	}

	for _, p := range notSupportedPlatforms {
		_, _, err := getSDKAttributes(p.os, p.arch)
		if err.Error() != fmt.Sprintf("unsupported platform: %s-%s", p.os, p.arch) {
			t.Errorf("unexpected error: %s", err)
		}
	}
}

func TestGetSDKDownloadURL(t *testing.T) {
	if runtime.GOOS == "darwin" && runtime.GOARCH == "amd64" {
		url, platform, _ := GetSDKDownloadURL()
		assert.Equal(t, url, "https://download.asperasoft.com/download/sw/sdk/transfer/1.1.1/osx-amd64-52a85ef.tar.gz")
		assert.Equal(t, platform, "osx-amd64")
	}
}
