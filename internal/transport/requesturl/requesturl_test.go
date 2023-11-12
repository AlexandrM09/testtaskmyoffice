package requesturl

import (
	"context"
	"net/http"
	"testing"
	model "github.com/AlexandrM09/testtaskmyoffice/internal/model"
	"github.com/stretchr/testify/assert"
)

type rusecase struct {
}

func (r *rusecase) Get(ctx context.Context, c http.Client, v model.Requestline) model.Requestline {
	v.Size = v.Size * 2
	return v
}
func TestRequesturltransport(t *testing.T) {
	rusecase := &rusecase{}
	requesturl := NewRequesturltransport(9, 1000, rusecase)
	ctx := context.Background()
	in,out:=requesturl.Run(ctx)

	input :=[]int64{1,2,3,4,5,6,7,8,9}
	var output [9] int64
	for i,v := range input {
		vrl := model.Requestline{
			Nline: int64(i),
			Size: int64(v),
		}
		in<-vrl
		output[i]=v*2
	}
	requesturl.Stop()
    for v:=range out{
		assert.Equal(t, output[v.Nline], v.Size, "equal size")
	}
}
