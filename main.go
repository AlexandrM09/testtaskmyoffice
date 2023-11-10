package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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
	countWorker     = 40
	msecondduration = 1000
)

func main() {
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
		go createrequestworker(msecondduration, &wg, in, out)
	}
	//output result
	go func(out chan requestline) {
		for v := range out {
			fmt.Printf("%d/n", v.Size)
		}
	}(out)
	//
	reader := bufio.NewReader(file)
	counter := int64(1)
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
		// fmt.Println(string(line))
	}
	close(in)
	wg.Wait()
	close(out)
	log.Println("Done")
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

//validate url
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
	defer resp.Body.Close()
	if err != nil {
		v.Err = fmt.Errorf("request get error:%w", err)
		return v
	}
	var buf bytes.Buffer
	v.Size, err = bodysize(io.TeeReader(resp.Body, &buf))
	if err != nil {
		v.Err = fmt.Errorf("body read error:%w", err)
	}
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
