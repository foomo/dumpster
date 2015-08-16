package handler

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/foomo/dumpster/dumpster"
	"github.com/julienschmidt/httprouter"
)

func SetupHandlers(r *httprouter.Router, d *dumpster.Dumpster) {
	handlerDump := &Dump{
		dumpster: d,
	}
	handlerDump.Register(r)
	handlerRestore := &Restore{
		dumpster: d,
	}
	handlerRestore.Register(r)
}

func jsonReply(data interface{}, w http.ResponseWriter) error {
	jsonBytes, err := json.MarshalIndent(data, "", "   ")
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
	log.Println("sent json reply:", len(jsonBytes))
	return nil
}

func errReply(w http.ResponseWriter, code int, err error) {
	log.Println("an error occurred", code, err.Error())
	w.WriteHeader(code)
	w.Write([]byte(err.Error()))
}

func extractJSONBodyIntoData(r *http.Request, data interface{}) error {
	jsonBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if len(jsonBytes) == 0 {
		return errors.New("body was empty")
	}
	err = json.Unmarshal(jsonBytes, &data)
	if err != nil {
		return err
	}
	return nil
}
