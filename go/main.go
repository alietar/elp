package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/alietar/elp/go/gpsfiles"
	"github.com/alietar/elp/go/server"
)

func main() {
	dlAll, dlSome := flagHandler()

	downloadDepartment(dlAll, dlSome)

	fmt.Println("\n\033[34m --- Starting the server --- \033[0m\n")
	server.Start()
}

func isThereDBFolder() bool {
	if _, err := os.Stat("./db"); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func flagHandler() (bool, bool) {
	downloadAllFlagPtr := flag.Bool("dl-all", false, "Downloads all the db")
	downloadSomeFlagPtr := flag.Bool("dl-some", false, "Downloads only the provided departement")
	forceDownloadFlagPtr := flag.Bool("dl-force", false, "Forces the download even if a db folder exists")

	flag.Parse()

	if *downloadAllFlagPtr && *downloadSomeFlagPtr {
		fmt.Println("You can't download all departements and some at the same time.")
		fmt.Println("Remove either the -dl-all flag or the -dl-some flag")
		os.Exit(3)
	}

	if (*downloadAllFlagPtr || *downloadSomeFlagPtr) && isThereDBFolder() && !*forceDownloadFlagPtr {
		fmt.Println("You want to download the DB but a db folder already exists,")
		fmt.Println("to overwrite the folder, use the -dl-force flag")
		os.Exit(3)
	}

	if !*downloadAllFlagPtr && !*downloadSomeFlagPtr && !isThereDBFolder() {
		fmt.Println("No DB folder found, please download the DB with the  -dl-all or -dl-some flag")
		os.Exit(3)
	}

	return *downloadAllFlagPtr, *downloadSomeFlagPtr
}

func downloadDepartment(dlAll, dlSome bool) {
	if dlAll {
		gpsfiles.DownloadAllDepartements()
	} else if dlSome {
		if len(flag.Args()) == 0 {
			fmt.Println("No department specified, please specify which department you want to download")
			fmt.Println("by adding their number at the end of the command")

			os.Exit(3)
		}

		var departementsNb []int

		for _, nbStr := range flag.Args() {
			nbInt, err := strconv.Atoi(nbStr)

			if err != nil {
				fmt.Println("Specified department number is not an int")
				os.Exit(3)
			}

			if nbInt < 1 || nbInt == 20 || nbInt > 95 {
				fmt.Println("Wrong department number")
				os.Exit(3)
			}

			departementsNb = append(departementsNb, nbInt)
		}

		for _, nb := range departementsNb {
			gpsfiles.DownloadUnzipDB(nb)
		}
	}
}
