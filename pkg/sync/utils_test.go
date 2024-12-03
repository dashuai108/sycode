package sync

import (
	"testing"
)

func TestExist(t *testing.T) {
	path := "/tmp/sycodeout/test1/aa.txt"
	exist, err := Exist(path)
	if err != nil {
		t.Error(err)
	}
	t.Logf("The path [%s] exist %t", path, exist)

}
