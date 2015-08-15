package handler

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func SetupHandlers(r *httprouter.Router) {
	d := &Dump{}
	d.Register(r)
}

func jsonReply(data interface{}, w http.ResponseWriter) error {
	jsonBytes, err := json.MarshalIndent(data, "", "   ")
	if err != nil {
		return err
	}
	w.Write(jsonBytes)
	return nil
}
