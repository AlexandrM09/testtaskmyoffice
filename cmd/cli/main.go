package main

import (
	"flag"
	"runtime"

	a "github.com/AlexandrM09/testtaskmyoffice/internal/app"
)

func main() {
	var flagconfig a.Appconfig
	// flags
	flag.StringVar(&flagconfig.Path, "path", "./source/testurl.txt", "path url file")
	flag.IntVar(&flagconfig.Cpucount, "cpucount", 6, "cpu count")
	flag.IntVar(&flagconfig.CountWorker, "countWorker", 10, "workers count")
	flag.IntVar(&flagconfig.Maxprocessurldurationmsec, "maxprocessurldurationmsec", 1000, "maximum duration get url request,msec")
	flag.IntVar(&flagconfig.Maxtotaldurationsecond, "maxtotaldurationsecond", 600, "maximum total duration,second")
	flag.Parse()
	// fmt.Printf("flags %v\n", flagconfig)
	//set GOMAXPROCS
	cpuCountSys := runtime.NumCPU()
	if flagconfig.Cpucount < cpuCountSys {
		runtime.GOMAXPROCS(flagconfig.Cpucount)
	}
	app := a.NewApp(flagconfig)
	app.Run()
}
