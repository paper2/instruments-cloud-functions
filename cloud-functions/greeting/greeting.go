package greeting

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func init() {
	functions.HTTP("Greeting", greeting)
}

// greetingHTTP is an HTTP Cloud Function.
func greeting(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hiya!")
	err := greetNext(context.Background())
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func greetNext(ctx context.Context) error {
	next := os.Getenv("NEXT_ENDPOINT")
	if next == "" {
		fmt.Println("I have no freinds :(")
		return nil
	}

	_, err := otelhttp.Get(ctx, next)
	if err != nil {
		return err
	}
	fmt.Println("I said hi to my friend :)")

	return nil
}
