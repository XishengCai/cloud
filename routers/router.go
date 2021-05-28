package routers

import (
	_ "cloud/docs"
	v1 "cloud/routers/api/v1"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

var (
	// httpHistogram prometheus 模型
	httpHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   "http_server",
		Subsystem:   "",
		Name:        "requests_seconds",
		Help:        "Histogram of response latency (seconds) of http handlers.",
		ConstLabels: nil,
		Buckets:     nil,
	}, []string{"method", "code", "uri"})
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	prometheus.MustRegister(httpHistogram)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	kubernetes := r.Group("/api/kubernetes/v1/")

	kubernetes.Use()
	{
		// cluster
		kubernetes.POST("/masters", v1.InstallKubernetes)
		kubernetes.POST("/slaves", v1.InstallKubernetesSlave)
		kubernetes.POST("/uninstall", nil)
		kubernetes.POST("/update", nil)
	}

	docker := r.Group("/api/docker/v1/")
	docker.Use()
	{
		docker.POST("/install", v1.InstallDocker)
		docker.POST("/uninstall", v1.UnInstallDocker)
	}

	return r
}
