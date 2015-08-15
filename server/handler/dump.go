package handler

import (
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
}

// List available dumps
func (d *Dump) GetDumps(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	dumps := make(map[string][]*responses.Dump)
	for dumpType, ddumps := range d.dumpster.GetDumps() {
		dumps[dumpType] = castDumps(ddumps)
	}
	jsonReply(dumps, w)
}

// List available dumps
func (d *Dump) GetDumpsForType(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
	dump, err := d.dumpster.GetDump(ps.ByName("type"), ps.ByName("id"))
	if err != nil {
		// 404
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
		return
	}
	http.ServeFile(w, r, dump.Filename)
}

func (d *Dump) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := d.dumpster.DeleteDump(ps.ByName("type"), ps.ByName("id"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte("sucessfully deleted"))
}

func (d *Dump) Create(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	createDumpRequest := &requests.CreateDump{}
	err := extractJSONBodyIntoData(r, createDumpRequest)
	if err != nil {
		panic(err)
	}
	dumpType := ps.ByName("type")
	metaData, err := d.dumpster.CreateDump(dumpType, createDumpRequest.ID, createDumpRequest.Comment)
	if err != nil {
		panic(err)
	}
	jsonReply(&responses.Dump{
		ID:      metaData.ID,
		Comment: metaData.Comment,
		Report:  metaData.Report,
		Errors:  metaData.Errors,
		Path:    getPath(metaData.Type, metaData.ID),
	}, w)
}

func (d *Dump) Restore(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func getPath(dumpType string, id string) string {
	return "/dump/" + url.QueryEscape(dumpType) + "/" + url.QueryEscape(id)
}
