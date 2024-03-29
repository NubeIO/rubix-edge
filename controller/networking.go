package controller

import (
	"github.com/NubeIO/lib-networking/networking"
	"github.com/NubeIO/rubix-edge/service/system"
	"github.com/gin-gonic/gin"
)

var nets = networking.New()

func (inst *Controller) Networking(c *gin.Context) {
	data, err := nets.GetNetworks()
	responseHandler(data, err, c)
}

func (inst *Controller) GetInterfacesNames(c *gin.Context) {
	data, err := nets.GetInterfacesNames()
	responseHandler(data, err, c)
}

func (inst *Controller) InternetIP(c *gin.Context) {
	data, err := nets.GetInternetIP()
	responseHandler(data, err, c)
}

func (inst *Controller) RestartNetworking(c *gin.Context) {
	data, err := inst.System.RestartNetworking()
	responseHandler(data, err, c)
}

func (inst *Controller) InterfaceUpDown(c *gin.Context) {
	var m system.NetworkingBody
	err := c.ShouldBindJSON(&m)
	if err != nil {
		responseHandler(nil, err, c)
		return
	}
	data, err := inst.System.InterfaceUpDown(m)
	responseHandler(data, err, c)
}

func (inst *Controller) InterfaceUp(c *gin.Context) {
	var m system.NetworkingBody
	err := c.ShouldBindJSON(&m)
	if err != nil {
		responseHandler(nil, err, c)
		return
	}
	data, err := inst.System.InterfaceUp(m)
	responseHandler(data, err, c)
}

func (inst *Controller) InterfaceDown(c *gin.Context) {
	var m system.NetworkingBody
	err := c.ShouldBindJSON(&m)
	if err != nil {
		responseHandler(nil, err, c)
		return
	}
	data, err := inst.System.InterfaceDown(m)
	responseHandler(data, err, c)
}
