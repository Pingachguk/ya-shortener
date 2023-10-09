package config

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Addr struct {
	Host string
	Port int
}

func (addr *Addr) Set(value string) error {
	hp := strings.Split(value, ":")
	if len(hp) != 2 {
		return errors.New("need address in a form host:port")
	}
	port, err := strconv.Atoi(hp[1])
	if err != nil {
		return err
	}
	addr.Host = hp[0]
	addr.Port = port
	return nil
}

func (addr *Addr) String() string {
	return fmt.Sprintf("%s:%d", addr.Host, addr.Port)
}

var defaultAddr Addr = Addr{
	Host: "0.0.0.0",
	Port: 8080,
}

var (
	AppAddr  = defaultAddr
	BaseAddr = defaultAddr
)

func InitConfig() {
	parseFlags()
}
