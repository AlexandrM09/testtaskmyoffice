package main

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"
)

type requestline struct {
	Nline    int64
	Url      string
	Readtime time.Duration
	Size     int64
	Err      error
}

const (
// cpucount                  = 4
// countWorker               = 10
// maxprocessurldurationmsec = 1000
// maxtotaldurationsecond    = 60 * 10 * time.Second
)

func main() {
	//to do
	//flags
	//log
	//gracefull shotdown
	//refactoring
	//make
	//test
	//github actions
	//pprof
	//README

	// flags
	cpucount := *flag.Int("cpucount", 2, "cpu count")
	countWorker := *flag.Int("countWorker", 10, "workers count")
	maxprocessurldurationmsec := *flag.Int("maxprocessurldurationmsec", 1000, "maximum duration get url request,msec")
	maxtotaldurationsecond := *flag.Int("maxtotaldurationsecond ", 60*10, "maximum total duration,second")
	//set NumCPU
	cpuCountSys := runtime.NumCPU()
	if cpucount >= cpuCountSys {
		runtime.GOMAXPROCS(cpuCountSys - 1)
	}
	if cpucount < cpuCountSys {
		runtime.GOMAXPROCS(cpucount)
	}
	start := time.Now()
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Duration(maxtotaldurationsecond)*time.Second)
	defer cancelFunc()
	defer func(start time.Time) {
		fmt.Printf("total processing time,msec: %d\n", time.Now().Sub(start).Milliseconds())
	}(start)
	// first open the file
	file, err := os.Open("./source/testurl.txt")
	if err != nil {
		log.Fatalf("could not open the file: %v", err)
	}
	defer file.Close()
	//---------
	in := make(chan requestline, countWorker)
	out := make(chan requestline, countWorker)
	var wg sync.WaitGroup
	for i := 0; i < countWorker; i++ {
		go createrequestworker(int64(maxprocessurldurationmsec), &wg, in, out)
	}
	//output result
	go func(out chan requestline) {
		for v := range out {
			if v.Err == nil {
				fmt.Printf("line source:%d,url:%s,content size,kB:%0.1f,processing time,msec: %d\n", v.Nline, v.Url, float64(v.Size)/1024.0, v.Readtime.Milliseconds())
			}
			if v.Err != nil {
				fmt.Printf("line source:%d,url:%s,error:%s\n", v.Nline, v.Url, v.Err)
			}
		}
	}(out)
	//
	reader := bufio.NewReader(file)
	counter := int64(1)
	var ctxresult string = "done"
l:
	for {
		line, err := read(reader)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("a real error happened here: %v\n", err)
		}
		v := requestline{
			Nline: counter,
			Url:   string(line),
		}
		in <- v
		counter++
		select {
		case <-ctx.Done():
			{
				ctxresult = "context timeout exceeded"
				break l
			}
		default:
		}
		// fmt.Println(string(line))
	}
	close(in)
	wg.Wait()
	close(out)
	log.Println(ctxresult)
}

// Read with Readline function

func read(r *bufio.Reader) ([]byte, error) {
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

// validate url
func validate(url string) error {
	return nil
}

// request worker
func createrequestworker(msecondduration int64, wg *sync.WaitGroup, in chan requestline, out chan requestline) {
	c := http.Client{Timeout: time.Duration(msecondduration) * time.Millisecond}
	wg.Add(1)
	for v := range in {
		out <- urlget(c, v)
	}
	wg.Done()
}

func urlget(c http.Client, v requestline) requestline {
	//vaidate url
	if err := validate(v.Url); err != nil {
		v.Err = fmt.Errorf("validate url error:%w", err)
		return v

	}
	//request
	start := time.Now()
	resp, err := c.Get(v.Url)
	v.Readtime = time.Now().Sub(start)
	// defer resp.Body.Close()
	if err != nil {
		v.Err = fmt.Errorf("request get error:%w", err)
		return v
	}
	var buf bytes.Buffer
	v.Size, err = bodysize(io.TeeReader(resp.Body, &buf))
	if err != nil {
		v.Err = fmt.Errorf("body read error:%w", err)
	}
	resp.Body.Close()
	return v
}
func bodysize(r io.Reader) (int64, error) {
	maxmem := 4096
	bytes := make([]byte, maxmem)
	var size int64
	for {
		read, err := r.Read(bytes)
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}
		size += int64(read)
	}
	return size, nil
}
