package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"runtime/pprof"
	"strconv"

	"github.com/alietar/elp/go/gpsfiles"
	"github.com/alietar/elp/go/server"
	"github.com/alietar/elp/go/tileutils"
)

func main() {
	dlAll, dlSome, accuracy, perfMode := flagHandler()

	if perfMode {
		for i := range 7 {
			path := fmt.Sprintf("perf/test_%d.prof", i)
			f, err := os.Create(path)
			if err != nil {
				log.Fatal(err)
			}

			pprof.StartCPUProfile(f)

			// tileutils.ComputeTiles(4.871928, 45.7838052, 0.3, gpsfiles.ACCURACY_1)
			nWorker := int(math.Pow(2, float64(i)))
			fmt.Printf("\n\nnWorker: %d\n", nWorker)
			for j := range 3 {
				fmt.Println(j)
				tileutils.ComputeTiles(4.979897, 45.784764, 2, gpsfiles.ACCURACY_1, nWorker) // Le grand large
				// tileutils.ComputeTiles(4.636917, 45.779077, 2, gpsfiles.ACCURACY_1, nWorker) // Dans la montagne
				// tileutils.ComputeTiles(4.871492, 45.763811, 3, gpsfiles.ACCURACY_1, nWorker) // En ville
			}

			pprof.StopCPUProfile()
		}
	} else {
		downloadDepartments(dlAll, dlSome, accuracy)

		fmt.Println("\n        _,--',   _._.--._____")
		fmt.Println(" .--.--';_'-.', \";_      _.,-'")
		fmt.Println(".'--'.  _.'    {`'-;_ .-.>.'")
		fmt.Println("      '-:_      )  / `' '=.")
		fmt.Println("        ) >     {_/,     /~)")
		fmt.Println("        |/               `^ .'")

		fmt.Println("\n\033[34m --- FIND REACHABLE ---\033[0m\n")
		fmt.Println("\n\033[34m -> Starting the server\033[0m\n")

		server.Start()
	}
}

func isThereDBFolder() bool {
	if _, err := os.Stat("./db"); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func flagHandler() (bool, bool, gpsfiles.MapAccuracy, bool) {
	downloadAllFlagPtr := flag.Bool("dl-all", false, "Downloads all the db")
	downloadSomeFlagPtr := flag.Bool("dl-some", false, "Downloads only the provided departement. Put them after all the flags!")
	accuracy1FlagPtr := flag.Bool("accuracy-1", false, "1M Accuracy")
	accuracy5FlagPtr := flag.Bool("accuracy-5", false, "5M Accuracy")
	accuracy25FlagPtr := flag.Bool("accuracy-25", false, "25M Accuracy")
	perfFlagPtr := flag.Bool("perf", false, "Doing perf tests")

	flag.Parse()

	if *perfFlagPtr {
		return false, false, gpsfiles.ACCURACY_25, true
	}

	if *downloadAllFlagPtr && *downloadSomeFlagPtr {
		fmt.Println("You can't download all departements and some at the same time.")
		fmt.Println("Remove either the -dl-all flag or the -dl-some flag")
		os.Exit(3)
	}

	if !*downloadAllFlagPtr && !*downloadSomeFlagPtr && !isThereDBFolder() {
		fmt.Println("No DB folder found, please download the DB with the -dl-all or -dl-some, and an accuracy flag")
		os.Exit(3)
	}

	var accuracy gpsfiles.MapAccuracy

	if *downloadAllFlagPtr || *downloadSomeFlagPtr {
		if !*accuracy1FlagPtr && !*accuracy5FlagPtr && !*accuracy25FlagPtr {
			fmt.Println("Please indicate at which accuracy you want to download")
			fmt.Println("the file with the flags -accuracy-1, -accuracy-5 or -accuracy-25")
			os.Exit(3)
		}

		if *downloadAllFlagPtr && *accuracy1FlagPtr {
			fmt.Println("Downloading at 1M the whole France is too much, use -dl-some")
			os.Exit(3)
		}

		if *accuracy1FlagPtr {
			accuracy = gpsfiles.ACCURACY_1
		} else if *accuracy5FlagPtr {
			accuracy = gpsfiles.ACCURACY_5
		} else if *accuracy25FlagPtr {
			accuracy = gpsfiles.ACCURACY_25
		} else {
			fmt.Println("Invalid accuracy, use either 1, 5 or 25")
			os.Exit(3)
		}
	}

	return *downloadAllFlagPtr, *downloadSomeFlagPtr, accuracy, false
}

func downloadDepartments(dlAll, dlSome bool, accuracy gpsfiles.MapAccuracy) {
	if dlAll {
		gpsfiles.DownloadAllDepartements(accuracy)

		fmt.Println("\n\033[32m\033[1mSuccessfully downloaded requested departments\033[0m")
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
			gpsfiles.DownloadUnzipDepartment(nb, accuracy)
		}

		fmt.Println("\n\033[32m\033[1mSuccessfully downloaded requested departments\033[0m")
	}
}
