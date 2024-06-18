package main

import (
	"pbench/cmd"

	_ "pbench/cmd/queryplan"
)

func main() {
	cmd.Execute()
}
