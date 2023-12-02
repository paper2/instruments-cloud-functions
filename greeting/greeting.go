package greeting

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func init() {
	// NOTE: Usually TraceProvider should call Shutdown() at the end of the program.
	//       It is mitigted by using ForceFlush() in InstrumentedHandler.
	tracerProvider := InitTracing()
	handler := InstrumentedHandler("greeting", greeting, tracerProvider)
	functions.HTTP("Greeting", handler)
}

// greetingHTTP is an HTTP Cloud Function.
func greeting(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hiya!")
	// sleep for extend span duration
	time.Sleep(100 * time.Millisecond)

	err := greetNext(r.Context())
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// greetNext calls the next function.
func greetNext(ctx context.Context) error {
	next := os.Getenv("NEXT_ENDPOINT")
	if next == "" {
		log.Println("I have no freinds :(")
		return nil
	}

	res, err := otelhttp.Get(ctx, next)
	if err != nil {
		return err
	}

	err = res.Body.Close()
	if err != nil {
		return err
	}

	log.Println("I said hi to my friend :)")

	return nil
}
