package common

import "testing"

func TestParserTemplate(t *testing.T) {
	type obj struct {
		Version string `json:"version"`
	}
	testCases := []struct {
		FilePath string
	}{
		{
			FilePath: "../../template/install_kubeadm.sh",
		},
	}

	for _, item := range testCases {
		buf, err := ParserTemplate(item.FilePath, obj{
			Version: "1.17.11",
		})

		if err != nil {
			t.Fatal(err)
		}

		t.Log(buf)

	}
}
