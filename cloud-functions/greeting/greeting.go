package greeting

import (
	"fmt"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("Greeting", greeting)
}

// greetingHTTP is an HTTP Cloud Function.
func greeting(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "Paper2"
	}
	fmt.Fprintf(w, "Hi, %s!", name)
}
