package unique

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateUniqueId(t *testing.T) {
	ids, _ := GenerateUniqueId(BizTypeWorker, 1)
	fmt.Println(ids)
	assert2 := assert.New(t)
	assert2.True(true)
}
