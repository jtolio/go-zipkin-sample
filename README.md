go-zipkin-sample
===========

This project is a toy example of how to use the
github.com/spacemonkeygo/monitor/trace Zipkin client library.

There's a random number service that just returns random numbers, and there's
a random message service that returns a random word from the list
"hello, world", using the random number service. Both services send Zipkin
spans to the collector service over UDP, which sends the spans to Zipkin
directly.

First, build the following binaries:
* github.com/jtolds/go-zipkin-sample/bin/collector
* github.com/jtolds/go-zipkin-sample/bin/rand
* github.com/jtolds/go-zipkin-sample/bin/ui

Then, after starting Zipkin locally (https://github.com/itszero/docker-zipkin
is helpful), run all three services (./rand & ./ui & ./collector &)

When you make a request to http://localhost:8079, a full Zipkin trace will
be sent to your Zipkin collector.
