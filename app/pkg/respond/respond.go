package respond

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  string `json:"result"`
}

func (res Response) ErrorHandle(w http.ResponseWriter, r *http.Request, code int, err error) {
	Respond(w, r, code, res)
}

func SuccessStrHandle(w http.ResponseWriter, r *http.Request, code int, err string) {
	Respond(w, r, code, Response{Status: "success", Message: err})
}

func Respond(w http.ResponseWriter, _ *http.Request, code int, data any) {
	w.WriteHeader(code)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
