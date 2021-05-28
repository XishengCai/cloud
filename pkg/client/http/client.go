package http

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"k8s.io/klog"
	"net/http"
)

func BuildHttpClient(caCrt, serverCert, serverKey []byte) (*http.Client, error) {
	if caCrt == nil || serverCert == nil || serverKey == nil {
		return nil, fmt.Errorf("证书不可以为空")
	}

	pool := x509.NewCertPool()

	if !pool.AppendCertsFromPEM(caCrt) {
		return nil, fmt.Errorf("pool.AppendCertsFromPEM failed")
	}

	cliCrt, err := tls.X509KeyPair(serverCert, serverKey)
	if err != nil {
		return nil, fmt.Errorf("parses a public/private, %v", err)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:      pool,
			Certificates: []tls.Certificate{cliCrt},
		},
	}
	return &http.Client{Transport: tr}, nil
}

func GetHttp(url string) (data []byte, err error) {
	klog.Infof("GetHttp url: %s", url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("charset", "UTF-8")

	clt := http.Client{}
	resp, err := clt.Do(req)
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	klog.V(4).Infof("resp: %+v", resp)

	return body, err
}
