package cli

import (
	"testing"
)

func TestImportapp(t *testing.T) {
	filepath := "/home/smallb/workspace/go/src/github.com/yansmallb/assigner/test/importapp.yaml"
	etcdpath := "http://192.168.11.51:2379"
	importapp(filepath, etcdpath)
}
