package controller

import (
	"fmt"
	"github.com/NubeIO/lib-systemctl-go/systemctl"
	"github.com/NubeIO/rubix-edge/model"
	"github.com/NubeIO/rubix-edge/service/system"
	"github.com/NubeIO/rubix-registry-go/rubixregistry"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"net/http"
)

type Controller struct {
	SystemCtl     *systemctl.SystemCtl
	RubixRegistry *rubixregistry.RubixRegistry
	System        *system.System
	FileMode      int
	Scheduler     *gocron.Scheduler
}

type Response struct {
	StatusCode   int         `json:"status_code"`
	ErrorMessage string      `json:"error_message"`
	Message      string      `json:"message"`
	Data         interface{} `json:"data"`
}

func responseHandler(body interface{}, err error, c *gin.Context, statusCode ...int) {
	var code int
	if err != nil {
		if len(statusCode) > 0 {
			code = statusCode[0]
		} else {
			code = http.StatusNotFound
		}
		msg := model.Message{
			Message: fmt.Sprintf("rubix-edge: %s", err.Error()),
		}
		c.JSON(code, msg)
	} else {
		if len(statusCode) > 0 {
			code = statusCode[0]
		} else {
			code = http.StatusOK
		}
		c.JSON(code, body)
	}
}
