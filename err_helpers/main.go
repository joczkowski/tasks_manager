package err_helpers

import "net/http"

func HandleWebErr(w http.ResponseWriter, err error, status int) {
	if err != nil {
		http.Error(w, err.Error(), status)
	}
}
