package http

import (
	"testing"
	"time"
)

func TestHttpGet(t *testing.T) {

	//url := "http://47.99.114.224:31560/agent/api/v1/cluster/overview?start=1592471772&end=1592471772"
	url := "http://47.99.114.224:31501/agent/api/v1/cluster/overview?start=1592472750&end=1592472750"
	url = "http://118.31.225.120"
	i := 1
	for {
		go func() {
			b, err := GetHttp(url)
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("time: %d, %s ", i, string(b))
		}()
		i++
		if i > 10000 {
			break
		}
		if i > 100 {
			time.Sleep(3 * time.Second)
		}

	}

	time.Sleep(1000 * time.Second)

}
