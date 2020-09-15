package valuate

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValuate(t *testing.T) {
	exp := "(val*60-(val*60)%250)>2000?(val*60-(val*60)%250):2000"
	mapVar := map[string]interface{}{"val": 2000}

	assert2 := assert.New(t)
	v, err := ValuateToFloat64(exp, mapVar)
	if err != nil {
		assert2.Fail("error", err)
		return
	}

	fmt.Println(v)
	assert2.True(v == 2000)
}
