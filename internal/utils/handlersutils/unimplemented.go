package handlersutils

import (
	"fmt"
	"net/http"
)

func Unimplemented(w http.ResponseWriter, r *http.Request) {
	apiError := NewAPIError(
		"UNIMPLEMENTED_ROUTE",
		fmt.Sprintf("Route '%s' is unimplemented. Method: %s", r.URL.Path, r.Method),
	)

	WriteResponseJSON(w, apiError, http.StatusAccepted)
}
