package main

import (
	"github.com/alietar/elp/go/algo"
	"github.com/alietar/elp/go/server"
)

func main() {
	algo.DownloadAllDepartements()

	server.Start()
}
