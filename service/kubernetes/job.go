package kubernetes

import (
	"encoding/json"
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"k8s.io/klog"
	"os"
	"os/signal"
)

var (
	RedisPool = &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", ":6379")
		},
	}
	installK8sQueue      = work.NewEnqueuer(InstallK8sJobNamespace, RedisPool)
	installK8sSlaveQueue = work.NewEnqueuer(InstallK8sSlaveJobNamespace, RedisPool)
)

const (
	installMaster = "install_master"
	installSlave  = "install_slave"

	InstallK8sJobNamespace      = "install_k8s_master"
	InstallK8sSlaveJobNamespace = "install_k8s_slave"
)

// ConvertJobArg convert to work.Q
func ConvertJobArg(i interface{}) (work.Q, error) {
	m := work.Q{}
	b, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &m)
	return m, err
}

// SetUpJob start job server
func SetUpJob() {
	go registerJob(InstallK8sJobNamespace, installMaster, InstallKuber{})
	go registerJob(InstallK8sSlaveJobNamespace, installSlave, InstallSlave{})
}

type job interface {
	ConsumeJob(j *work.Job) error
	Log(job *work.Job, next work.NextMiddlewareFunc) error
	Export(job *work.Job) error
}

type Context struct{}

func registerJob(namespace, jobName string, j job) {
	pool := work.NewWorkerPool(j, 10, namespace, RedisPool)

	pool.Job(jobName, j.ConsumeJob)
	pool.Middleware(j.Log)
	pool.JobWithOptions("export", work.JobOptions{Priority: 10, MaxFails: 1}, j.Export)
	pool.Start()

	klog.Info("register job ", jobName)
	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan
	// Stop the pool
	pool.Stop()
}
