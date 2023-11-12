package requesturl

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	model "github.com/AlexandrM09/testtaskmyoffice/internal/model"
)

type Requesturlusecase struct {
}

func Newrequesturlusecase() *Requesturlusecase {
	return &Requesturlusecase{}
}
func (r *Requesturlusecase) Get(ctx context.Context, c http.Client, v model.Requestline) model.Requestline {
	//vaidate url
	if err := validate(v.Url); err != nil {
		v.Err = fmt.Errorf("validate url error:%w", err)
		return v
	}
	//request
	start := time.Now()
	resp, err := c.Get(v.Url)
	v.Readtime = time.Since(start)
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

// bodysize
func bodysize(r io.Reader) (int64, error) {
	maxmem := 4096
	bytes := make([]byte, maxmem)
	var size int64
	for {
		read, err := r.Read(bytes)
		size += int64(read)
		if err == io.EOF {
			// size += int64(read)
			break
		}
		if err != nil {
			return 0, err
		}
	}
	return size, nil
}

// validate url
func validate(urls string) error {
	_, err := url.ParseRequestURI(urls)
	return err
}
