package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHostname(t *testing.T) {
	name := Hostname()
	fmt.Println(name)
	assert2 := assert.New(t)
	assert2.True(true)
}

func TestLocalIP(t *testing.T) {
	name := LocalIP()
	fmt.Println(name)
	assert2 := assert.New(t)
	assert2.True(true)
}

func TestOSInfo(t *testing.T) {
	fmt.Println(OSInfo())
	assert2 := assert.New(t)
	assert2.True(true)
}