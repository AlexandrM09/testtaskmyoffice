package main

import (
	"flag"
	"runtime"

	a "github.com/AlexandrM09/testtaskmyoffice/internal/app"
)

func main() {
	var flagconfig a.Appconfig
	// flags
	flagconfig.Path = *flag.String("path", "./source/testurl.txt", "path url file")
	flagconfig.Cpucount = *flag.Int("cpucount", 6, "cpu count")
	flagconfig.CountWorker = *flag.Int("countWorker", 10, "workers count")
	flagconfig.Maxprocessurldurationmsec = *flag.Int("maxprocessurldurationmsec", 1000, "maximum duration get url request,msec")
	flagconfig.Maxtotaldurationsecond = *flag.Int("maxtotaldurationsecond ", 60*10, "maximum total duration,second")
	//set GOMAXPROCS
	cpuCountSys := runtime.NumCPU()
	if flagconfig.Cpucount < cpuCountSys {
		runtime.GOMAXPROCS(flagconfig.Cpucount)
	}
	app := a.NewApp(flagconfig)
	app.Run()
}
