package main

import (
	"fmt"

	"github.com/elsni/lagerator/args"
	"github.com/elsni/lagerator/data"
	"github.com/elsni/lagerator/loggi"
)

// main is the program entry point.
func main() {
	//ui.TestForm()
	fmt.Println()
	data.Db.Load()
	args.ProcessArgs()
	fmt.Println()
	loggi.Log.Print(false)
}
