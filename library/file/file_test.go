package file

import (
	"fmt"
	"testing"
)

func TestReadLineAt(t *testing.T) {
	v, err := ReadLineAt("dev", 2)
	if err != nil {
		fmt.Println(err)
		t.Fail()
		return
	}

	if v != "xxxxxxx\nxxxxxxx" {
		t.Fail()
	}
}
