package o7tracker

import (
	"encoding/json"
	"net/http"
)

// ErrorToJSON outputs nice JSON errors to ResponseWritter
func ErrorToJSON(w http.ResponseWriter, err error) {
	if err != nil {
		json, _ := json.Marshal(map[string]interface{}{
			"status": "error",
			"msg":    err.Error(),
		})

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(json)
		return
	}
}
