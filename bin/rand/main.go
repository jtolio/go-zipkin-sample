package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"code.google.com/p/go.net/context"
	"github.com/spacemonkeygo/monitor/trace"
	"github.com/spacemonkeygo/monitor/trace/gen-go/zipkin"
)

type handler struct{}

func (h *handler) ServeHTTP(ctx context.Context, w http.ResponseWriter,
	r *http.Request) {
	defer trace.Trace(&ctx)(nil)
	time.Sleep(5 * time.Millisecond)
	fmt.Fprint(w, h.RandomNumber(ctx))
}

func (h *handler) RandomNumber(ctx context.Context) (num int) {
	defer trace.Trace(&ctx)(nil)

	time.Sleep(3 * time.Millisecond)
	return rand.Int()
}

func RandomNumberService(address string) error {
	return http.ListenAndServe(address, trace.ContextWrapper(trace.TraceHandler(
		&handler{})))
}

func main() {
	c, err := trace.NewUDPCollector("127.0.0.1:8082", 128)
	if err != nil {
		panic(err)
	}
	trace.Configure(1, true, &zipkin.Endpoint{
		Ipv4:        127*256*256*256 + 1,
		Port:        8081,
		ServiceName: "randhost",
	})
	trace.RegisterTraceCollector(c)

	panic(RandomNumberService("127.0.0.1:8081"))
}
