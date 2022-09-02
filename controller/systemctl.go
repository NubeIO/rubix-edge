package controller

import (
	"github.com/NubeIO/lib-rubix-installer/installer"
	"github.com/gin-gonic/gin"
)

// SystemCtlAction start, stop, enable, disable a service
func (inst *Controller) SystemCtlAction(c *gin.Context) {
	var m *installer.SystemCtlBody
	err := c.ShouldBindJSON(&m)
	if err != nil {
		responseHandler(nil, err, c)
		return
	}
	data, err := inst.EdgeApp.App.SystemCtlAction(m)
	responseHandler(data, err, c)
}

func (inst *Controller) SystemCtlStatus(c *gin.Context) {
	var m *installer.SystemCtlBody
	err := c.ShouldBindJSON(&m)
	if err != nil {
		responseHandler(nil, err, c)
		return
	}
	data, err := inst.EdgeApp.App.SystemCtlStatus(m)
	responseHandler(data, err, c)
}

// ServiceMassAction start, stop, enable, disable a service
func (inst *Controller) ServiceMassAction(c *gin.Context) {
	var m *installer.SystemCtlBody
	err := c.ShouldBindJSON(&m)
	if err != nil {
		responseHandler(nil, err, c)
		return
	}
	data, err := inst.EdgeApp.App.ServiceMassAction(m)
	responseHandler(data, err, c)
}

func (inst *Controller) ServiceMassStatus(c *gin.Context) {
	var m *installer.SystemCtlBody
	err := c.ShouldBindJSON(&m)
	if err != nil {
		responseHandler(nil, err, c)
		return
	}
	data, err := inst.EdgeApp.App.ServiceMassStatus(m)
	responseHandler(data, err, c)
}
