package apps

import (
	"context"
	"fmt"
	"github.com/NubeIO/git/pkg/git"
	"github.com/NubeIO/lib-systemctl-go/builder"
	"github.com/NubeIO/lib-systemctl-go/ctl"
	"github.com/NubeIO/lib-systemctl-go/systemctl"
	log "github.com/sirupsen/logrus"
	pprint "gthub.com/NubeIO/rubix-cli-app/pkg/helpers/print"
	"gthub.com/NubeIO/rubix-cli-app/service/apps/app"
	"time"
)

var err error

func New(inst *Apps) (*Apps, error) {

	a := &app.App{}
	selection := &app.Selection{
		AppName: inst.AppName,
		Version: "latest",
	}
	selectApp, err := a.SelectApp(selection)
	if err != nil {
		return nil, err
	}
	selectApp = selectApp.GetApp()
	pprint.PrintJOSN(selectApp)
	opts := &git.AssetOptions{
		Owner: selectApp.Owner,
		Repo:  selectApp.Repo,
		Tag:   selectApp.Version,
		Arch:  selectApp.Arch,
	}
	ctx := context.Background()
	inst.gitClient = git.NewClient(inst.Token, opts, ctx)
	return inst, err
}

//type RespDownload struct {
//	AssetName string `json:"asset_name"`
//	GitError  string `json:"git_error"`
//}

type RespBuilder struct {
	BuilderErr string `json:"builder_err"`
}

type RespInstall struct {
	InstallResp *ctl.InstallResp `json:"install_resp"`
}

type Apps struct {
	Token              string `json:"git_token"`
	Version            string `json:"tag"`
	AppName            string
	DownloadPath       string `json:"download_path"`         //home/user
	DownloadPathSubDir string `json:"download_path_sub_dir"` //home/user /bios
	gitClient          *git.Client
}

func (inst *Apps) GitDownload(destination string) (*git.DownloadResponse, error) {
	return inst.gitClient.Download(destination)
}

func (inst *Apps) GenerateServiceFile() (*RespBuilder, error) {
	ret := &RespBuilder{}
	newService := "nubeio-rubix-bios"
	description := "BIOS comes with default OS, non-upgradable"
	user := "root"
	directory := "/data/rubix-bios-app"
	execCmd := "/data/rubix-bios-app/rubix-bios -p 1615 -g /data/rubix-bios -d data -c config -a apps --prod --auth  --device-type amd64 --token 1234"

	bld := &builder.SystemDBuilder{
		Description:      description,
		User:             user,
		WorkingDirectory: directory,
		ExecStart:        execCmd,
		SyslogIdentifier: "rubix-bios",
		WriteFile: builder.WriteFile{
			Write:    true,
			FileName: newService,
			Path:     "/data",
		},
	}

	err = bld.Build()
	if err != nil {
		ret.BuilderErr = err.Error()
		return ret, err
	}
	return ret, nil
}

//InstallService a new linux service
//	- service: the service name (eg: rubix-bios)
//	- path: the service file path and name (eg: "/tmp/rubix-bios.service")
func (inst *Apps) InstallService(service, path string) (*RespInstall, error) {
	ret := &RespInstall{}
	//path := "/tmp/nubeio-rubix-bios.service"
	timeOut := 30
	ser := ctl.New(service, path)
	opts := systemctl.Options{Timeout: timeOut}
	installOpts := ctl.InstallOpts{
		Options: opts,
	}
	ser.InstallOpts = installOpts
	ret.InstallResp = ser.Install()
	fmt.Println("full install error", err)
	if err != nil {
		fmt.Println("full install error", err)
	}

	time.Sleep(8 * time.Second)

	status, err := systemctl.Status(service, systemctl.Options{})
	if err != nil {
		log.Errorf("service found: %s: %v", service, err)
	}
	fmt.Println(status)
	return ret, nil
}
