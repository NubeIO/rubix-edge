package rubix

import (
	"fmt"
	fileutils "github.com/NubeIO/lib-dirs/dirs"
	"testing"
)

func TestApp_uploadApp(t *testing.T) {
	var err error
	homeDir, _ := fileutils.Dir()
	fmt.Println(homeDir)
	app := New(&App{DataDir: "/data", Perm: nonRoot, HostDownloadPath: fmt.Sprintf("%s/Downloads", homeDir)})

	appName := "rubix-wires"
	appInstallName := "wires-builds"
	serviceName := "nubeio-rubix-wires"
	appVersion := "v2.7.2"
	appZip := "wires-builds-2.7.2.zip"

	fmt.Println(appName, appInstallName, serviceName, appVersion, appZip)

	err = app.uploadApp(appName, appInstallName, appVersion, nil)

	fmt.Println(err)
}
