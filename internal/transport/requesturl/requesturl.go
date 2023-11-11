package requesturl

import (
	"context"
	"net/http"
	"sync"
	"time"

	model "github.com/AlexandrM09/testtaskmyoffice/internal/model"
)

type Requesturlusecase interface {
	Get(ctx context.Context, c http.Client, v model.Requestline) model.Requestline
}
type Requesturltransport struct {
	requesturl                Requesturlusecase
	in                        chan model.Requestline
	out                       chan model.Requestline
	countWorker               int
	maxprocessurldurationmsec int
	wg                        *sync.WaitGroup
}

func NewRequesturltransport(countWorker int, maxprocessurldurationmsec int, requesturl Requesturlusecase) *Requesturltransport {

	return &Requesturltransport{
		requesturl:                requesturl,
		in:                        make(chan model.Requestline, countWorker),
		out:                       make(chan model.Requestline, countWorker),
		countWorker:               countWorker,
		maxprocessurldurationmsec: maxprocessurldurationmsec,
		wg:                        &sync.WaitGroup{},
	}
}
func (r *Requesturltransport) Run(ctx context.Context) (in,out chan model.Requestline) {
	for i := 0; i < r.countWorker; i++ {
		go func() {
			c := http.Client{Timeout: time.Duration(r.maxprocessurldurationmsec) * time.Millisecond}
			r.wg.Add(1)
			for v := range r.in {
				r.out <- r.requesturl.Get(ctx, c, v)
			}
			r.wg.Done()
		}()
	}
	return r.in,r.out
}
func (r *Requesturltransport) Stop() {
	close(r.in)
	r.wg.Wait()
	close(r.out)
}
