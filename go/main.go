package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"runtime/pprof"
	"strconv"

	"github.com/alietar/elp/go/gpsfiles"
	"github.com/alietar/elp/go/server"
	"github.com/alietar/elp/go/tileutils"
)

func main() {
	dlAll, dlSome, accuracy, perfMode, port, decompressPath := flagHandler()

	if perfMode {
		nCPU := runtime.NumCPU()

		nExploreWorker := int(math.Sqrt(float64(nCPU)))
		nFileWorker := int(nCPU / nExploreWorker)

		fmt.Println(nExploreWorker)
		fmt.Println(nFileWorker)

		pathCpu := fmt.Sprintf("perf/cpu_%d_worker.prof", nExploreWorker)

		fCpu, err := os.Create(pathCpu)
		if err != nil {
			log.Fatal(err)
		}

		pprof.StartCPUProfile(fCpu)

		for j := range 20 {
			fmt.Println(j)
			tileutils.ComputeTiles(4.979897, 45.784764, 2, gpsfiles.ACCURACY_1, nExploreWorker, nFileWorker) // Le grand large
		}

		pprof.StopCPUProfile()

		// 	nWorker := int(math.Pow(2, float64(i)))

		// 	pathCpu := fmt.Sprintf("perf/cpu_%d_worker.prof", nWorker)
		// 	pathMem := fmt.Sprintf("perf/mem_%d_worker.prof", nWorker)
		// 	fCpu, err := os.Create(pathCpu)
		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}
		// 	fMem, err := os.Create(pathMem)
		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}

		// 	pprof.StartCPUProfile(fCpu)

		// 	fmt.Printf("\n\nnWorker: %d\n", nWorker)
		// 	for j := range 3 {
		// 		fmt.Println(j)
		// 		tileutils.ComputeTiles(4.979897, 45.784764, 2, gpsfiles.ACCURACY_1, nWorker) // Le grand large
		// 		tileutils.ComputeTiles(4.636917, 45.779077, 2, gpsfiles.ACCURACY_1, nWorker) // Dans la montagne
		// 		tileutils.ComputeTiles(4.871492, 45.763811, 3, gpsfiles.ACCURACY_1, nWorker) // En ville
		// 	}

		// 	pprof.StopCPUProfile()
		// 	pprof.WriteHeapProfile(fMem)

		// 	fmt.Println("\nFinished profiling")
		// }
	} else {
		downloadDepartments(dlAll, dlSome, accuracy, decompressPath)

		fmt.Println("\n        _,--',   _._.--._____")
		fmt.Println(" .--.--';_'-.', \";_      _.,-'")
		fmt.Println(".'--'.  _.'    {`'-;_ .-.>.'")
		fmt.Println("      '-:_      )  / `' '=.")
		fmt.Println("        ) >     {_/,     /~)")
		fmt.Println("        |/               `^ .'")

		fmt.Println("\n\033[34m --- FIND REACHABLE ---\033[0m\n")
		fmt.Printf("\n\n\033[34m -> Starting the server on port %d\033[0m\n", port)

		server.Start(port)
	}
}

func isThereDBFolder() bool {
	if _, err := os.Stat("./db"); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func flagHandler() (bool, bool, gpsfiles.MapAccuracy, bool, int, string) {
	downloadAllFlagPtr := flag.Bool("dl-all", false, "Downloads all the departement")
	downloadSomeFlagPtr := flag.Bool("dl-some", false, "Downloads only the provided departement. Put them after all the flags!")
	accuracyFlagPtr := flag.Int("accuracy", 25, "Specify accuracy in meters")
	decompressFlagPtr := flag.String("7z", "7z", "Path to 7z")
	perfFlagPtr := flag.Bool("perf", false, "Doing perf tests")

	portFlagPtr := flag.Int("port", 8080, "HTTP Server's port")

	flag.Parse()

	if *perfFlagPtr {
		return false, false, gpsfiles.ACCURACY_25, true, 0, ""
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
		switch *accuracyFlagPtr {
		case 1:
			accuracy = gpsfiles.ACCURACY_1
		case 5:
			accuracy = gpsfiles.ACCURACY_5
		case 25:
			accuracy = gpsfiles.ACCURACY_25
		default:
			fmt.Println("Invalid accuracy, use either 1, 5 or 25")
			os.Exit(3)
		}

		if *downloadAllFlagPtr && accuracy == gpsfiles.ACCURACY_1 {
			fmt.Println("Downloading at 1M the whole France is too much, use -dl-some")
			os.Exit(3)
		}
	}

	return *downloadAllFlagPtr, *downloadSomeFlagPtr, accuracy, false, *portFlagPtr, *decompressFlagPtr
}

func downloadDepartments(dlAll, dlSome bool, accuracy gpsfiles.MapAccuracy, decompressPath string) {
	if dlAll {
		gpsfiles.DownloadAllDepartements(accuracy, decompressPath)

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
			gpsfiles.DownloadUnzipDepartment(nb, accuracy, decompressPath)
		}

		fmt.Println("\n\033[32m\033[1mSuccessfully downloaded requested departments\033[0m")
	}
}
