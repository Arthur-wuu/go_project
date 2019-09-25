package utils

import (
	"context"
	"fmt"
	"github.com/satori/go.uuid"
	"sync"
)

// Trace and monitor request
// TODO: not finish, not useful

type ApiRequest struct {
	VersionApi string // api
	Id         string // uuid
}

func (ar ApiRequest) BuildUuid() error {
	uuid, err := uuid.NewV4()
	if err != nil {
		ar.Id = uuid.String()
	}

	return err
}

type ApiRequestTracer struct {
	stopChan chan bool

	apiChan chan ApiRequest

	apiFinChan chan string

	apiRequestMap map[string]ApiRequest
}

func NewApiRequestTracer() *ApiRequestTracer {
	art := &ApiRequestTracer{}
	art.stopChan = make(chan bool)
	art.apiChan = make(chan ApiRequest, 1000)
	art.apiFinChan = make(chan string, 1000)
	art.apiRequestMap = make(map[string]ApiRequest)
	return art
}

func (art *ApiRequestTracer) StartTracer(ctx context.Context, wg *sync.WaitGroup) {
	go func(ctx2 context.Context, wg2 *sync.WaitGroup, art2 *ApiRequestTracer) {
		wg2.Add(1)
		defer wg2.Done()

		for {
			select {
			case <-ctx2.Done():
				art2.stopChan <- true
				return
			case ar := <-art2.apiChan:
				art2.apiRequestMap[ar.Id] = ar
				fmt.Println("register api:", ar)
			case id := <-art2.apiFinChan:
				delete(art2.apiRequestMap, id)
				fmt.Println("finish api:", id)
			}
		}
	}(ctx, wg, art)

	<-art.stopChan
	fmt.Println("I am quit Api request tracer")
}

func (art *ApiRequestTracer) TraceApi(versionApi string) string {
	uuid, err := uuid.NewV4()
	if err != nil {
		return ""
	}

	ar := ApiRequest{VersionApi: versionApi, Id: uuid.String()}

	art.apiChan <- ar
	return ar.Id
}

func (art *ApiRequestTracer) FinishTraceApi(id string) {
	art.apiFinChan <- id
	return
}
