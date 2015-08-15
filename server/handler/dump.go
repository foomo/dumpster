package handler

import (
	"fmt"
	"net/http"

	"github.com/foomo/dumpster/server/vo/responses"
	"github.com/julienschmidt/httprouter"
)

// Dump handler
type Dump struct {
}

func (d *Dump) Register(r *httprouter.Router) {
	r.GET("/dump/:name", d.List)
	r.GET("/dump/:name/:id", d.Get)
	r.DELETE("/dump/:name/Id", d.Delete)
}

// List available dumps
func (d *Dump) List(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	jsonReply(map[string][]*responses.Dump{
		"foo": []*responses.Dump{
			&responses.Dump{
				ID:   "foo",
				Path: "/dump/foo/a",
			},
		},
	}, w)
}

func (d *Dump) Get(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Write([]byte(fmt.Sprintln(ps)))
}

func (d *Dump) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func (d *Dump) Create(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func (d *Dump) Restore(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}
