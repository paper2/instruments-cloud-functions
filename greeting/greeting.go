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
	// NOTE: Usually TracerProvider should call Shutdown() at the end of the program.
	//       It is difficult to do so in Cloud Functions.
	//       This issue can be mitigated by using ForceFlush() to flush spans.
	tracerProvider := InitTracing()
	handler := InstrumentedHandler("greeting", greeting, tracerProvider)
	functions.HTTP("Greeting", handler)
}

// greeting is the function's core logic.
// It resoponses "Hiya!" and calls the next function if NEXT_ENDPOINT is set.
func greeting(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hiya!")

	// sleep for extending the span's duration
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

	// call the next function.
	// otelhttp sends a trace context to the next function.
	res, err := otelhttp.Get(ctx, next)
	if err != nil {
		return err
	}

	// Must close the response body to avoid leaking connections.
	err = res.Body.Close()
	if err != nil {
		return err
	}

	log.Println("I said hi to my friend :)")

	return nil
}
