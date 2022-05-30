package controller

import (
	"github.com/gin-gonic/gin"
	"gthub.com/NubeIO/rubix-cli-app/service/apps"
)

func getAppStoreBody(ctx *gin.Context) (dto *apps.Store, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func (inst *Controller) GetAppStores(c *gin.Context) {
	data, err := inst.DB.GetAppStores()
	if err != nil {
		reposeHandler(nil, err, c)
		return
	}
	reposeHandler(data, err, c)
}

func (inst *Controller) GetAppStore(c *gin.Context) {
	data, err := inst.DB.GetAppStore(c.Params.ByName("uuid"))
	if err != nil {
		reposeHandler(nil, err, c)
		return
	}
	reposeHandler(data, err, c)
}

func (inst *Controller) CreateAppStore(c *gin.Context) {
	var m *apps.Store
	err = c.ShouldBindJSON(&m)
	data, err := inst.DB.CreateAppStore(m)
	if err != nil {
		reposeHandler(nil, err, c)
		return
	}
	reposeHandler(data, err, c)
}

func (inst *Controller) UpdateAppStore(c *gin.Context) {
	body, _ := getAppStoreBody(c)
	data, err := inst.DB.UpdateAppStore(c.Params.ByName("uuid"), body)
	if err != nil {
		reposeHandler(nil, err, c)
		return
	}
	reposeHandler(data, err, c)
}

func (inst *Controller) DeleteAppStore(c *gin.Context) {
	q, err := inst.DB.DeleteAppStore(c.Params.ByName("uuid"))
	if err != nil {
		reposeHandler(nil, err, c)
	} else {
		reposeHandler(q, err, c)
	}
}

func (inst *Controller) DropAppStores(c *gin.Context) {
	data, err := inst.DB.DropAppStores()
	if err != nil {
		reposeHandler(nil, err, c)
		return
	}
	reposeHandler(data, err, c)
}
