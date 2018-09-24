package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/transport/zipkin"
	zipkinC "github.com/uber/jaeger-client-go/zipkin"
)

var (
	zipkinURL = flag.String("url",
		"http://zipkin.istio-system.svc.cluster.local:9411/api/v1/spans", "Zipkin server URL")
	serverPort = flag.String("port", "8000", "server port")
	traceLabel = "opentracing-go-nethttp-demo"
)

func getTime(w http.ResponseWriter, r *http.Request) {
	log.Print("Received getTime request")
	t := time.Now()
	ts := t.Format("Mon Jan _2 15:04:05 2006")
	io.WriteString(w, fmt.Sprintf("The time is %s", ts))
}

func redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/gettime", 301)
}

func main() {
	flag.Parse()

	zipkinPropagator := zipkinC.NewZipkinB3HTTPHeaderPropagator()
	injector := jaeger.TracerOptions.Injector(opentracing.HTTPHeaders, zipkinPropagator)
	extractor := jaeger.TracerOptions.Extractor(opentracing.HTTPHeaders, zipkinPropagator)
	zipkinSharedRPCSpan := jaeger.TracerOptions.ZipkinSharedRPCSpan(true)

	// Set Transport type
	transport, err := zipkin.NewHTTPTransport(
		*zipkinURL,
		zipkin.HTTPBatchSize(1),
		zipkin.HTTPLogger(jaeger.StdLogger),
	)
	if err != nil {
		log.Fatalf("HTTP transport error: %v", err)
	}

	// Create tracer
	tracer, closer := jaeger.NewTracer(
		traceLabel,
		jaeger.NewConstSampler(true),
		jaeger.NewRemoteReporter(transport),
		injector,
		extractor,
		zipkinSharedRPCSpan,
	)

	// Create HTTP Server
	http.HandleFunc("/gettime", getTime)
	http.HandleFunc("/", redirect)
	log.Printf("Starting server on port %s", *serverPort)
	err = http.ListenAndServe(
		fmt.Sprintf(":%s", *serverPort),
		// use nethttp.Middleware to enable OpenTracing for server
		nethttp.Middleware(tracer, http.DefaultServeMux))
	if err != nil {
		log.Fatalf("Cannot start server: %s", err)
	}

	closer.Close()
}
