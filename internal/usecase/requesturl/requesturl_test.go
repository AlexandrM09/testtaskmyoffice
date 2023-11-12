package requesturl

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	model "github.com/AlexandrM09/testtaskmyoffice/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestRequesturlusecase(t *testing.T) {
	requesturl := Newrequesturlusecase()
	input := "expected"
	s := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Send response to be tested
		_, err := rw.Write([]byte(input))
		if err != nil {
			t.Fail()
		}
	}))
	defer s.Close()

	c := http.Client{Timeout: 1 * time.Second}
	v := model.Requestline{
		Nline: 1,
		Url:   s.URL,
	}
	ctx := context.Background()
	actual := requesturl.Get(ctx, c, v)
	assert.NoError(t, actual.Err, "error")
	fmt.Printf("body size: %d\n", actual.Size)
	assert.Equal(t, int64(len(input)), actual.Size, "equal size body")

}

func TestInvalidURL(t *testing.T) {
	requesturl := Newrequesturlusecase()
	c := http.Client{Timeout: 50 * time.Millisecond}
	v := model.Requestline{
		Nline: 1,
		Url:   "ht&@-tp://:aa",
	}
	ctx := context.Background()
	actual := requesturl.Get(ctx, c, v)
	assert.Error(t, actual.Err)
	assert.Equal(t, int64(0), actual.Size, "equal size body")
}
func TestTimeout(t *testing.T) {
	requesturl := Newrequesturlusecase()
	input := "expected"
	s := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Send response to be tested
		time.Sleep(50 * time.Millisecond)
		_, err := rw.Write([]byte(input))
		if err != nil {
			t.Fail()
		}
	}))
	defer s.Close()

	c := http.Client{Timeout: 30 * time.Millisecond}
	v := model.Requestline{
		Nline: 1,
		Url:   s.URL,
	}
	ctx := context.Background()
	actual := requesturl.Get(ctx, c, v)
	assert.Error(t, actual.Err)
	assert.Equal(t, int64(0), actual.Size, "equal size body")
}
func TestBodyReadError(t *testing.T) {
	requesturl := Newrequesturlusecase()
	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Length", "1")
		}),
	)
	defer s.Close()

	c := http.Client{Timeout: 30 * time.Millisecond}
	v := model.Requestline{
		Nline: 1,
		Url:   s.URL,
	}
	ctx := context.Background()
	actual := requesturl.Get(ctx, c, v)
	assert.Error(t, actual.Err)
	assert.Equal(t, int64(0), actual.Size, "equal size body")
}
