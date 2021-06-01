package kubernetes

import (
	"encoding/json"
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"k8s.io/klog"
)

var (
	redisPool = &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", ":6379")
		},
	}
	installK8sQueue      = work.NewEnqueuer(installK8sJobNamespace, redisPool)
	installK8sSlaveQueue = work.NewEnqueuer(installK8sSlaveJobNamespace, redisPool)
)

const (
	installMaster = "install_master"
	installSlave  = "install_slave"

	installK8sJobNamespace      = "install_k8s_master"
	installK8sSlaveJobNamespace = "install_k8s_slave"
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
	go registerJob(installK8sJobNamespace, installMaster, InstallKuber{})
	go registerJob(installK8sSlaveJobNamespace, installSlave, InstallSlave{})
}

type job interface {
	ConsumeJob(j *work.Job) error
	Log(job *work.Job, next work.NextMiddlewareFunc) error
	Export(job *work.Job) error
}

func registerJob(namespace, jobName string, j job) {
	pool := work.NewWorkerPool(j, 10, namespace, redisPool)

	pool.Job(jobName, j.ConsumeJob)
	pool.Middleware(j.Log)
	pool.JobWithOptions("export", work.JobOptions{Priority: 10, MaxFails: 1}, j.Export)
	pool.Start()

	klog.Info("register job ", jobName)
	// Wait for a signal to quit:
	//signalChan := make(chan os.Signal, 1)
	//signal.Notify(signalChan, os.Interrupt, os.Kill)
	//<-signalChan

	// Stop the pool
	//pool.Stop()
}
