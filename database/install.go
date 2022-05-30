package dbase

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gthub.com/NubeIO/rubix-cli-app/service/apps"
)

type App struct {
	AppName string `json:"app_name"`
	Version string `json:"version"`
	Token   string `json:"token"`
}

type InstallResponse struct {
	Message    string     `json:"message"`
	Error      string     `json:"error"`
	InstallLog InstallLog `json:"log"`
}

type InstallLog struct {
	GetAppFromStore string `json:"get_app_from_store"`
	AppInstall      string `json:"-"`
	MakeDownload    string `json:"make_download"`
	GitDownload     string `json:"git_download"`
	MakeInstallDir  string `json:"make_install_dir"`
	UnpackBuild     string `json:"unpack_build"`
	GenerateService string `json:"generate_service"`
	InstallService  string `json:"install_service"`
	CleanUp         string `json:"clean_up"`
}

// ok messages
const (
	selectAppStore    = "ok"
	makeDownload      = "ok"
	gitDownload       = "ok"
	makeNewApp        = "installed a new app"
	makeInstallDir    = "ok"
	unpackBuild       = "ok"
	generateService   = "ok"
	installService    = "ok"
	cleanUp           = "ok"
	updateExistingApp = ""
)

// not ok messages
const (
	selectAppStoreErr    = "this app is was not found in the app store, try flow-framework, rubix-wires"
	makeDownloadErr      = "issue on trying to make the path to download the zip folder"
	gitDownloadErr       = "error on git download"
	makeNewAppErr        = "failed to make a new app"
	makeInstallDirErr    = "unable to make the install dir for the app"
	unpackBuildErr       = "unable to unzip the build"
	generateServiceErr   = "unable to make the app service file"
	installServiceErr    = "unable to install the app"
	cleanUpErr           = "unable to clean up the install"
	updateExistingAppErr = ""
)

func (db *DB) InstallApp(body *App) (*InstallResponse, error) {
	resp := &InstallResponse{}
	app, err := db.installApp(body)
	if err != nil {
		resp.Error = err.Error()
		return app, err
	}
	resp.InstallLog = app.InstallLog
	resp.Error = "no errors"
	resp.Message = fmt.Sprintf("install ok! %s", app.InstallLog.AppInstall)
	return resp, err
}

func (db *DB) installApp(body *App) (*InstallResponse, error) {
	resp := &InstallResponse{}

	appStore, err := db.GetAppStoreByName(body.AppName)
	if err != nil {
		resp.InstallLog.GetAppFromStore = selectAppStoreErr
		return resp, err
	}
	resp.InstallLog.GetAppFromStore = selectAppStore
	installedApp := &apps.App{
		AppStoreName:     appStore.Name,
		AppStoreUUID:     appStore.UUID,
		InstalledVersion: body.Version,
	}

	var inst = &apps.Apps{
		Token:   body.Token,
		Perm:    apps.Permission,
		Version: body.Version,
		App:     appStore,
	}
	newApp, err := apps.New(inst)
	if err != nil {
		log.Errorln("new app: failed to init a new app", err)
		return resp, err
	}

	if err = inst.MakeDownloadDir(); err != nil {
		resp.InstallLog.MakeDownload = makeDownloadErr
		return resp, err
	}
	resp.InstallLog.MakeDownload = makeDownload
	download, err := newApp.GitDownload(inst.App.DownloadPath)
	if err != nil {
		log.Errorf("git: download error %s \n", err.Error())
		resp.InstallLog.GitDownload = gitDownloadErr
		return resp, err
	}
	assetTag := download.RepositoryRelease.GetTagName()
	resp.InstallLog.GitDownload = fmt.Sprintf("installed version: %s", assetTag)
	if err = inst.MakeInstallDir(); err != nil {
		resp.InstallLog.MakeInstallDir = makeInstallDirErr
		return resp, err
	}
	resp.InstallLog.MakeInstallDir = makeInstallDir
	if err = inst.UnpackBuild(); err != nil {
		resp.InstallLog.UnpackBuild = unpackBuildErr
		return resp, err
	}
	resp.InstallLog.UnpackBuild = unpackBuild
	tmpFileDir := newApp.App.DownloadPath
	if _, err = newApp.GenerateServiceFile(newApp, tmpFileDir); err != nil {
		log.Errorf("make service file build: failed error:%s \n", err.Error())
		resp.InstallLog.GenerateService = generateServiceErr
		return resp, err
	}
	resp.InstallLog.GenerateService = generateService
	tmpServiceFile := fmt.Sprintf("%s/%s.service", tmpFileDir, newApp.App.ServiceName)
	if _, err = newApp.InstallService(newApp.App.ServiceName, tmpServiceFile); err != nil {
		resp.InstallLog.InstallService = installServiceErr
		return resp, err
	}
	resp.InstallLog.InstallService = installService
	if err = inst.CleanUp(); err != nil {
		resp.InstallLog.CleanUp = cleanUpErr
		return resp, err
	}
	resp.InstallLog.CleanUp = cleanUp
	installedApp.InstalledVersion = assetTag
	app, existingApp, err := db.AddApp(installedApp)
	if err != nil {
		resp.InstallLog.AppInstall = makeNewAppErr
		return resp, err
	}
	if existingApp { // if it was existing app update the version
		app.InstalledVersion = assetTag
		_, err := db.UpdateApp(app.UUID, app)
		if err != nil {
			resp.InstallLog.AppInstall = fmt.Sprintf("an existing app was installed error:%s", err.Error())
			return resp, err
		}
		resp.InstallLog.AppInstall = fmt.Sprintf("an existing app was installed upgraded from: %s to: %s", app.InstalledVersion, assetTag)
	} else {
		resp.InstallLog.AppInstall = makeNewApp
	}

	log.Infof(fmt.Sprintf("an existing app was installed upgraded from:%s to:%s", app.InstalledVersion, assetTag))

	return resp, err

}
