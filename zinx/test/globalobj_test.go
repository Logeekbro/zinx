package test

import (
	"testing"
	"zinx/utils"
)

func TestGlobalObj(t *testing.T) {
	utils.GlobalObject.LoadConfig()
}
