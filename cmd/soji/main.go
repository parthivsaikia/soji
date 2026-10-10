package main

import (
	"os"

	"github.com/parthivsaikia/soji/internal/cli"
)

func main() {
	os.Exit(cli.Execute())
}
