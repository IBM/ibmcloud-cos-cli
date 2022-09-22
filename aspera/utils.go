package aspera

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultAddress = "127.0.0.1"
	defaultPort    = "55002"
)

func SDKDir() string {
	if sdk_path, ok := os.LookupEnv("ASPERA_SDK_PATH"); ok {
		return sdk_path
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return path.Join(home, ".aspera_sdk")
}

func TransferdBinPath() string {
	daemonName := "asperatransferd"
	if runtime.GOOS == "windows" {
		daemonName = "asperatransferd.exe"
	}
	return path.Join(SDKDir(), "bin", daemonName)
}

func GetSDKDownloadURL() (url string, platform string, err error) {
	platform, ext, err := getSDKAttributes(runtime.GOOS, runtime.GOARCH)
	if err != nil {
		return
	}
	url = fmt.Sprintf("%s/%s/%s-%s.%s", prefix, version, platform, commit, ext)
	return
}

func getSDKAttributes(os, arch string) (platform string, ext string, err error) {
	platforms := map[string][]string{
		"darwin":  {"amd64"},
		"linux":   {"amd64", "ppc64le", "s390x"},
		"windows": {"amd64"},
		"aix":     {"ppc64"},
	}

	ext = "tar.gz"

	if archs, ok := platforms[os]; ok {
		for _, a := range archs {
			if a == arch {
				if os == "darwin" {
					os = "osx"
				}
				if os == "windows" {
					ext = "zip"
				}
				return fmt.Sprintf("%s-%s", os, arch), ext, nil
			}
		}
	}
	return "", "", fmt.Errorf("unsupported platform: %s-%s", os, arch)

}

func defaultConnection() (cc *grpc.ClientConn, err error) {
	optInsecure := grpc.WithTransportCredentials(insecure.NewCredentials())
	target := fmt.Sprintf("%s:%s", defaultAddress, defaultPort)
	cc, err = grpc.Dial(target, optInsecure)
	return
}
