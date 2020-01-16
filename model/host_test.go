package model

import "testing"

func TestGetHost(t *testing.T){
	h, err := GetHost("47.56.114.4")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", h)
}