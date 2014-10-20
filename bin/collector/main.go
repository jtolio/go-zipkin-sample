package main

import (
	"github.com/spacemonkeygo/monitor/trace"
	"github.com/spacemonkeygo/spacelog"
)

func main() {
	spacelog.MustSetup("collector", spacelog.SetupConfig{
		Level:  "info",
		Format: "{{ColorizeLevel .Level}}{{.Message}}{{.Reset}}"})

	collector, err := trace.NewScribeCollector("127.0.0.1:9410")
	if err != nil {
		panic(err)
	}
	defer collector.Close()

	panic(trace.RedirectPackets("127.0.0.1:8082", collector))
}
