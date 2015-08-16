package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/foomo/dumpster/dumpster"
	"github.com/foomo/dumpster/server/vo/requests"
	"github.com/foomo/dumpster/server/vo/responses"
	"github.com/julienschmidt/httprouter"
)

// Dump handler
type Dump struct {
	dumpster *dumpster.Dumpster
}

func NewDump(dumpster *dumpster.Dumpster) *Dump {
	return &Dump{
		dumpster: dumpster,
	}
}

func (d *Dump) Register(r *httprouter.Router) {
	r.GET("/dump", d.GetDumps)
	r.GET("/dump/:type", d.GetDumpsForType)
	r.GET("/dump/:type/:id", d.Get)
	r.DELETE("/dump/:type/:id", d.Delete)
	r.POST("/dump/:type", d.Create)
	r.GET("/dumpremote/:name", d.GetRemoteDumps)
}

func (d *Dump) GetRemoteDumps(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Println("serving remote dumps for", ps.ByName("name"))
	remote, ok := d.dumpster.Remotes[ps.ByName("name")]
	if !ok {
		errReply(w, http.StatusNotFound, errors.New("unknown remote"))
		return
	}
	response, err := http.Get(remote.EndPoint + "/dump")
	if err != nil {
		errReply(w, http.StatusBadGateway, err)
		return
	}
	defer response.Body.Close()
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		errReply(w, http.StatusInternalServerError, err)
	}
	dumps := make(map[string][]*responses.Dump)
	filteredDumps := make(map[string][]*responses.Dump)
	err = json.Unmarshal(responseBytes, &dumps)
	for dumpType, remoteDumps := range dumps {
		w.Write([]byte(fmt.Sprintln(dumpType, remoteDumps)))
		dumpConfig, ok := d.dumpster.Dumps[dumpType]
		if ok && len(dumpConfig.Restore.Program) > 0 {
			filteredDumps[dumpType] = remoteDumps
		}
	}
	jsonReply(filteredDumps, w)
}

// List available dumps
func (d *Dump) GetDumps(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Println("serving all dumps")
	dumps := make(map[string][]*responses.Dump)
	for dumpType, ddumps := range d.dumpster.GetDumps() {
		dumps[dumpType] = castDumps(ddumps)
	}
	jsonReply(dumps, w)
}

// List available dumps
func (d *Dump) GetDumpsForType(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Println("serving dumps for type", ps.ByName("type"))
	jsonReply(d.dumpster.GetDumpsForType(ps.ByName("type")), w)
}

func castDumps(ddumps []*dumpster.Dump) []*responses.Dump {
	dumps := []*responses.Dump{}
	for _, ddump := range ddumps {
		dumps = append(dumps, &responses.Dump{
			ID:      ddump.ID,
			Created: ddump.Created,
			Comment: ddump.Comment,
			Report:  ddump.Report,
			Errors:  ddump.Errors,
			Path:    getPath(ddump.Type, ddump.ID),
		})
	}
	return dumps
}

func (d *Dump) Get(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Println("serving dump", ps.ByName("type"), ps.ByName("id"))
	dump, err := d.dumpster.GetDump(ps.ByName("type"), ps.ByName("id"))
	if err != nil {
		code := http.StatusInternalServerError
		if err == dumpster.ErrNotFound {
			code = http.StatusNotFound
		}
		errReply(w, code, err)
		return
	}
	http.ServeFile(w, r, dump.Filename)
}

func (d *Dump) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Println("deleting dump", ps.ByName("type"), ps.ByName("id"))
	err := d.dumpster.DeleteDump(ps.ByName("type"), ps.ByName("id"))
	if err != nil {
		errReply(w, http.StatusInternalServerError, err)
		return
	}
	w.Write([]byte("sucessfully deleted"))
}

func (d *Dump) Create(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	dumpType := ps.ByName("type")
	log.Println("creating dump in", dumpType)
	createDumpRequest := &requests.CreateDump{}
	err := extractJSONBodyIntoData(r, createDumpRequest)
	if err != nil {
		errReply(w, http.StatusBadRequest, err)
		return
	}

	metaData, err := d.dumpster.CreateDump(dumpType, createDumpRequest.ID, createDumpRequest.Comment)
	if err != nil {
		errReply(w, http.StatusBadRequest, err)
		return
	}
	jsonReply(&responses.Dump{
		ID:      metaData.ID,
		Comment: metaData.Comment,
		Report:  metaData.Report,
		Errors:  metaData.Errors,
		Path:    getPath(metaData.Type, metaData.ID),
	}, w)
}

func getPath(dumpType string, id string) string {
	return "/dump/" + url.QueryEscape(dumpType) + "/" + url.QueryEscape(id)
}
