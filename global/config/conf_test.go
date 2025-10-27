package config

import (
	"fmt"
	"testing"
)

func TestConInit(t *testing.T) {
	ConfInit()
	fmt.Print(Conf.Jwt)
	fmt.Print(Conf.Service)
}
