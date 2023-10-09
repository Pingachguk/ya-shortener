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

type SchemeAddr struct {
	Addr
	Scheme string
}

func (addr *Addr) Set(value string) error {
	hp := strings.Split(value, ":")
	if len(hp) == 2 {
		port, err := strconv.Atoi(hp[1])
		if err != nil {
			return err
		}
		addr.Port = port
	}

	addr.Host = hp[0]
	return nil
}

func (addr *Addr) String() string {
	if addr.Port != 0 {
		return fmt.Sprintf("%s:%d", addr.Host, addr.Port)
	} else {
		return addr.Host
	}
}

func (addr *SchemeAddr) Set(value string) error {
	shp := strings.Split(value, "://")
	if len(shp) != 2 {
		return errors.New("need address in a form scheme://host:port")
	}
	addr.Addr.Set(shp[1])
	addr.Scheme = shp[0]
	return nil
}

func (addr *SchemeAddr) String() string {
	return fmt.Sprintf("%s://%s", addr.Scheme, addr.Addr.String())
}

var (
	defaultAddr Addr = Addr{
		Host: "0.0.0.0",
		Port: 8080,
	}
	schemeDefaultAddr SchemeAddr = SchemeAddr{
		defaultAddr,
		"http",
	}
)

var (
	AppAddr  = defaultAddr
	BaseAddr = schemeDefaultAddr
)

func InitConfig() {
	parseFlags()
}
