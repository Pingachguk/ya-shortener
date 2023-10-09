package config

import (
	"flag"
)

func parseFlags() {
	flag.Var(&AppAddr, "a", "Application address")
	flag.Var(&BaseAddr, "b", "Base address")

	flag.Parse()
}
