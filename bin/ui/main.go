package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/spacemonkeygo/errors"

	"golang.org/x/net/context"
	"gopkg.in/spacemonkeygo/monitor.v1/trace"
	"gopkg.in/spacemonkeygo/monitor.v1/trace/gen-go/zipkin"
)

var (
	FakeError = errors.NewClass("Fake Error")
)

type handler struct {
	messages   []string
	random_url string
}

func (h *handler) ServeHTTP(ctx context.Context, w http.ResponseWriter,
	r *http.Request) {
	defer trace.Trace(&ctx)(nil)
	time.Sleep(4 * time.Millisecond)
	msg, err := h.RandomMessage(ctx)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Fprintln(w, msg)
}

func (h *handler) RandomMessage(ctx context.Context) (msg string, err error) {
	defer trace.Trace(&ctx)(&err)
	time.Sleep(6 * time.Millisecond)

	for i := 0; i < 5; i++ {
		go h.fakeError(ctx)
	}

	num, err := h.randomNumber(ctx)
	if err != nil {
		return "", err
	}

	return h.getMessage(ctx, num)
}

func (h *handler) randomNumber(ctx context.Context) (num int, err error) {
	defer trace.Trace(&ctx)(&err)

	time.Sleep(7 * time.Millisecond)

	req, err := http.NewRequest("GET", h.random_url, nil)
	if err != nil {
		return 0, err
	}

	resp, err := trace.TraceRequest(ctx, http.DefaultClient, req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(string(data))
}

func (h *handler) fakeError(ctx context.Context) (err error) {
	defer trace.Trace(&ctx)(&err)
	time.Sleep(3 * time.Millisecond)
	return FakeError.New("test error")
}

func (h *handler) getMessage(ctx context.Context, num int) (msg string,
	err error) {
	defer trace.Trace(&ctx)(&err)
	time.Sleep(8 * time.Millisecond)
	return h.messages[num%len(h.messages)], nil
}

func UIService(address string, messages []string, random_url string) error {
	return http.ListenAndServe(address, trace.ContextWrapper(trace.TraceHandler(
		&handler{messages: messages, random_url: random_url})))
}

func main() {
	c, err := trace.NewUDPCollector("127.0.0.1:8082", 128)
	if err != nil {
		panic(err)
	}
	trace.Configure(1, true, &zipkin.Endpoint{
		Ipv4:        127*256*256*256 + 1,
		Port:        8079,
		ServiceName: "uihost"})
	trace.RegisterTraceCollector(c)

	panic(UIService("127.0.0.1:8079", []string{"hello", "world"},
		"http://127.0.0.1:8081/random/number"))
}
