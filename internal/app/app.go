package App

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"time"

	model "github.com/AlexandrM09/testtaskmyoffice/internal/model"
	requrltransport "github.com/AlexandrM09/testtaskmyoffice/internal/transport/requesturl"
	requrlusecase "github.com/AlexandrM09/testtaskmyoffice/internal/usecase/requesturl"

	//"github.com/golangci/golangci-lint/pkg/result"
	"golang.org/x/exp/slog"
)

type Appconfig struct {
	Path                      string
	Cpucount                  int
	CountWorker               int
	Maxprocessurldurationmsec int
	Maxtotaldurationsecond    int
}

type requesturltransport interface {
	Run(ctx context.Context) (in, out chan model.Requestline)
	Stop()
}
type App struct {
	logger           *slog.Logger
	flagconfig       Appconfig
	requesturl       requesturltransport
	successfullcount int
	errorcount       int
	reasonexit       string
}

func NewApp(flagconfig Appconfig) *App {
	//log
	jsonHandler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{})
	requesturl := requrltransport.NewRequesturltransport(flagconfig.CountWorker,
		flagconfig.Maxprocessurldurationmsec, requrlusecase.Newrequesturlusecase())
	return &App{
		logger:     slog.New(jsonHandler),
		flagconfig: flagconfig,
		requesturl: requesturl,
	}
}

func (a *App) Run() {
	//to do
	//make
	//test
	//github actions
	//pprof
	//README
	start := time.Now()
	//os signal handler
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)
	//total timeout
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Duration(a.flagconfig.Maxtotaldurationsecond)*time.Second)
	defer cancelFunc()
	//result string
	defer func(start time.Time) {
		fmt.Printf("total processing time,msec: %d,count of successfully processed urls: %d,count of urls processed with errors: %d,reason exit: %s\n", time.Since(start).Milliseconds(), a.successfullcount, a.errorcount, a.reasonexit)
	}(start)
	// open source the file
	file, err := os.Open(a.flagconfig.Path)
	if err != nil {
		log.Fatalf("could not open the file: %v", err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	//start processing
	in, out := a.requesturl.Run(ctx)
	//output result
	go func(out chan model.Requestline) {
		for v := range out {
			fmt.Print(a.selectoutputstr(v))
		}
	}(out)
	//read file and send url
	a.reasonexit = inputurl(ctx, quit, a.logger, reader, in)
	a.requesturl.Stop()
}

// select output result
func (a *App) selectoutputstr(v model.Requestline) string {
	if v.Err == nil {
		a.successfullcount++
		return fmt.Sprintf("line source:%d,url:%s,content size,kB:%0.1f,processing time,msec: %d\n", v.Nline, v.Url, float64(v.Size)/1024.0, v.Readtime.Milliseconds())
	}
	a.errorcount++
	a.logger.Error(v.Err.Error())
	return fmt.Sprintf("line source:%d,url:%s,error:%s\n", v.Nline, v.Url, v.Err)

}

// read file and send url
func inputurl(ctx context.Context, ossignal chan os.Signal, logger *slog.Logger, reader *bufio.Reader, in chan model.Requestline) string {
	counter := int64(1)
	for {
		select {
		case <-ossignal:
			{
				return "os signal quit"
			}
		case <-ctx.Done():
			{
				return "context timeout exceeded"
			}
		default:
			{
				line, err := readline(reader)
				if err != nil {
					if err == io.EOF {
						return "source file eof"

					}
					logger.Error("error read file: %v\n", err)
					return "error read file"

				}
				v := model.Requestline{
					Nline: counter,
					Url:   string(line),
				}
				in <- v
				counter++
			}
		}
	}
}

// long line read
func readline(r *bufio.Reader) ([]byte, error) {
	var (
		isPrefix = true
		err      error
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return ln, err
}
