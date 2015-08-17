package handler

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/foomo/dumpster/dumpster"
	"github.com/foomo/dumpster/server/vo/requests"
	"github.com/foomo/dumpster/server/vo/responses"
	"github.com/julienschmidt/httprouter"
)

var (
	ErrRemoteBadGateway = errors.New("can not reach remote server")
	ErrUnknownRemote    = errors.New("unknown remote")
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
	r.GET("/dumpremote", d.GetAllRemoteDumps)
	r.GET("/dumpremote/:name", d.GetRemoteDumps)
}

func (d *Dump) GetAllRemoteDumps(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	allDumps := []*responses.Dump{}
	for remoteName := range d.dumpster.Remotes {
		dumps, err := d.getRemoteDumps(remoteName)
		if err == nil {
			allDumps = append(allDumps, dumps...)
		}
	}
	jsonReply(allDumps, w)
}

func (d *Dump) GetRemoteDumps(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Println("serving remote dumps for")
	remoteName := ps.ByName("name")
	dumps, err := d.getRemoteDumps(remoteName)
	if err != nil {
		code := http.StatusInternalServerError
		switch true {
		case err == ErrUnknownRemote:
			code = http.StatusNotFound
		case err == ErrRemoteBadGateway:
			code = http.StatusBadGateway
		}
		errReply(w, code, err)
		return
	}
	jsonReply(dumps, w)
}

func (d *Dump) getRemoteDumps(remoteName string) (dumps []*responses.Dump, err error) {
	remote, ok := d.dumpster.Remotes[remoteName]
	if !ok {
		return nil, ErrUnknownRemote
	}
	response, err := http.Get(remote.EndPoint + "/dump")
	if err != nil {
		return nil, ErrRemoteBadGateway
	}
	defer response.Body.Close()
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	dumps = []*responses.Dump{}
	filteredDumps := []*responses.Dump{}
	err = json.Unmarshal(responseBytes, &dumps)
	for _, remoteDump := range dumps {
		dumpConfig, ok := d.dumpster.Dumps[remoteDump.DumpType]
		if ok && len(dumpConfig.Restore.Program) > 0 {
			remoteDump.Remote = remoteName
			filteredDumps = append(filteredDumps, remoteDump)
		}
	}
	return filteredDumps, nil
}

// List available dumps
func (d *Dump) GetDumps(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Println("serving all dumps")
	dumps := []*responses.Dump{}
	for _, ddumps := range d.dumpster.GetDumps() {
		dumps = append(dumps, castDumps(ddumps)...)
	}
	jsonReply(dumps, w)
}

// List available dumps
func (d *Dump) GetDumpsForType(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Println("serving dumps for type", ps.ByName("type"))
	dumps, err := d.dumpster.GetDumpsForType(ps.ByName("type"))
	if err != nil {
		errReply(w, http.StatusNotFound, err)
	}
	jsonReply(dumps, w)
}

func castDumps(ddumps []*dumpster.Dump) []*responses.Dump {
	dumps := []*responses.Dump{}
	for _, ddump := range ddumps {
		dumps = append(dumps, &responses.Dump{
			ID:       ddump.ID,
			Created:  ddump.Created,
			DumpType: ddump.Type,
			Comment:  ddump.Comment,
			Report:   ddump.Report,
			Errors:   ddump.Errors,
			Path:     getPath(ddump.Type, ddump.ID),
		})
	}
	return dumps
}

func (d *Dump) Get(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Println("serving dump", ps.ByName("type"), ps.ByName("id"))
	dump, err := d.dumpster.GetDump(ps.ByName("type"), ps.ByName("id"))
	if err != nil {
		code := http.StatusInternalServerError
		switch true {
		case err == dumpster.ErrDumpNotFound:
		case err == dumpster.ErrDumpTypeNotFound:
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
	textReply(w, http.StatusOK, "successfully deleted")
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
		ID:       metaData.ID,
		Comment:  metaData.Comment,
		DumpType: metaData.Type,
		Report:   metaData.Report,
		Errors:   metaData.Errors,
		Path:     getPath(metaData.Type, metaData.ID),
	}, w)
}

func getPath(dumpType string, id string) string {
	return "/dump/" + url.QueryEscape(dumpType) + "/" + url.QueryEscape(id)
}
