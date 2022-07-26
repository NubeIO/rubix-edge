package router

import (
	"fmt"
	"github.com/NubeIO/lib-rubix-installer/installer"
	"github.com/NubeIO/rubix-edge/controller"
	dbase "github.com/NubeIO/rubix-edge/database"
	"github.com/NubeIO/rubix-edge/pkg/config"
	"github.com/NubeIO/rubix-edge/pkg/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"io"
	"os"
	"time"
)

func Setup(db *gorm.DB) *gin.Engine {
	engine := gin.New()

	// Set gin access logs
	fileLocation := fmt.Sprintf("%s/edge.access.log", config.Config.GetAbsDataDir())
	f, err := os.OpenFile(fileLocation, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		logger.Logger.Errorf("Failed to create access log file: %v", err)
	} else {
		gin.SetMode(viper.GetString("gin.loglevel"))
		gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	}

	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
	engine.Use(cors.New(cors.Config{
		AllowMethods: []string{"GET", "POST", "DELETE", "OPTIONS", "PUT", "PATCH"},
		AllowHeaders: []string{
			"X-FLOW-Key", "Authorization", "Content-Type", "Upgrade", "Origin",
			"Connection", "Accept-Encoding", "Accept-Language", "Host",
		},
		ExposeHeaders:          []string{"Content-Length"},
		AllowCredentials:       true,
		AllowAllOrigins:        true,
		AllowBrowserExtensions: true,
		MaxAge:                 12 * time.Hour,
	}))

	appDB := &dbase.DB{
		DB: db,
	}

	rubixApps := installer.New(&installer.App{})

	api := controller.Controller{DB: appDB, Rubix: rubixApps}

	apiRoutes := engine.Group("/api")

	apps := apiRoutes.Group("/apps")
	{
		apps.POST("/add", api.AddUploadApp)
		apps.POST("/service/upload", api.UploadService)
		apps.POST("/service/install", api.InstallService)
	}
	appControl := apiRoutes.Group("/apps/control")
	{
		appControl.POST("/", api.AppService)
		appControl.POST("/bulk", api.AppService)
	}

	device := apiRoutes.Group("/device")
	{
		device.GET("/", api.GetDeviceInfo)
		device.POST("/", api.AddDeviceInfo)
		device.PATCH("/", api.UpdateDeviceInfo)
	}

	system := apiRoutes.Group("/system")
	{
		system.GET("/ping", api.Ping)
		system.GET("/time", api.HostTime)
		system.GET("/product", api.GetProduct)
		system.POST("/scanner", api.RunScanner)
	}

	networking := apiRoutes.Group("/networking")
	{
		networking.GET("/networks", api.Networking)
		networking.GET("/interfaces", api.GetInterfacesNames)
		networking.GET("/internet", api.InternetIP)
		networking.GET("/update/schema", api.GetIpSchema)
		networking.POST("/update/dhcp", api.SetDHCP)
		networking.POST("/update/static", api.SetStaticIP)
	}

	files := apiRoutes.Group("/files")
	{
		files.GET("/walk", api.WalkFile)  // /api/files/walk?file=/data
		files.GET("/list", api.ListFiles) // /api/files/list?file=/data
		files.POST("/rename", api.RenameFile)
		files.POST("/copy", api.CopyFile)
		files.POST("/move", api.MoveFile)
		files.POST("/upload", api.UploadFile)
		files.POST("/download", api.DownloadFile)
		files.DELETE("/delete", api.DeleteFile)
	}

	dirs := apiRoutes.Group("/dirs")
	{
		dirs.POST("/create", api.CreateDir)
		dirs.POST("/copy", api.CopyDir)
		dirs.DELETE("/delete", api.DeleteDir)
	}

	zip := apiRoutes.Group("/zip")
	{
		zip.POST("/unzip", api.Unzip)
		zip.POST("/zip", api.ZipDir)
	}

	return engine
}
