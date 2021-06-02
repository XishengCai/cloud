package v1

import (
	"cloud/pkg/app"
	"cloud/service/kubernetes"
	"net/http"

	"github.com/gin-gonic/gin"
	"k8s.io/klog"
)

// InstallKubernetes godoc
// @Summary install kubernets
// @Description install kubernetes master
// @Tags install kubernetes master
// @Accept  json
// @Produce  json
// @Param cluster body models.InstallKubernetesParameter true "install kubernetes master"
// @Failure 201 {object} app.Response
// @Router /kubernetes/v1/masters [post]
func InstallKubernetes(ctx *gin.Context) {
	k := kubernetes.InstallKuber{}
	err := ctx.BindJSON(&k)
	if err != nil {
		app.HandleError(ctx, http.StatusBadRequest, err)
		return
	}
	klog.Infof("InstallKubernetes parameter: %+v", k)
	app.HandleError(ctx, http.StatusCreated, k.Install())
}

// InstallKubernetesSlave godoc
// @Summary install kubernetes slave
// @Description install kubernetes slave
// @Tags install kubernetes slave
// @Accept  json
// @Produce  json
// @Param cluster body models.InstallSlaveParam true "install kubernetes slave"
// @Failure 201 {object} app.Response
// @Router /kubernetes/v1/slaves [post]
func InstallKubernetesSlave(ctx *gin.Context) {
	k := kubernetes.InstallSlave{}
	err := ctx.ShouldBindJSON(&k)
	if err != nil {
		app.HandleError(ctx, http.StatusBadRequest, err)
		return
	}
	klog.Infof("InstallKubernetes nodes parameter: %+v", k)
	app.HandleDataAndError(ctx,http.StatusCreated,k.Install())
}
