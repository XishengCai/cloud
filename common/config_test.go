package common

import (
	"testing"
)

func TestTomlInit(t *testing.T) {
	t.Logf("%+v", tomlConf)
	a := Base64Encode("caicai12")
	b, err := Base64Decode(a)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("a: %s, b: %s", a, b)

}
