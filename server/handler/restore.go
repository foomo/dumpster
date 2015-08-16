package handler

import (
	"log"
	"net/http"

	"github.com/foomo/dumpster/dumpster"
	"github.com/foomo/dumpster/server/vo/responses"
	"github.com/julienschmidt/httprouter"
)

// Dump handler
type Restore struct {
	dumpster *dumpster.Dumpster
}

func NewRestore(dumpster *dumpster.Dumpster) *Restore {
	return &Restore{
		dumpster: dumpster,
	}
}

func (restore *Restore) Register(r *httprouter.Router) {
	r.POST("/restore/:type/:id", restore.Restore)
	r.POST("/restoreremote/:remote/:type/:id", restore.RestoreRemote)
}

// Restore a dump
func (restore *Restore) Restore(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Println("restoring dump", ps.ByName("type"), ps.ByName("id"))
	restoreReport, restoreErrors, err := restore.dumpster.RestoreDump(ps.ByName("type"), ps.ByName("id"))
	restoreReply(w, restoreReport, restoreErrors, err)
}

// Restore a dump
func (restore *Restore) RestoreRemote(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Println("restoring remote dump", ps.ByName("remote"), ps.ByName("type"), ps.ByName("id"))
	restoreReport, restoreErrors, err := restore.dumpster.RestoreRemoteDump(ps.ByName("remote"), ps.ByName("type"), ps.ByName("id"))
	restoreReply(w, restoreReport, restoreErrors, err)
}

func restoreReply(w http.ResponseWriter, restoreReport string, restoreErrors string, err error) {
	if err != nil {
		errReply(w, http.StatusInternalServerError, err)
		return
	}
	jsonReply(&responses.RestoreReport{
		Report: restoreReport,
		Errors: restoreErrors,
	}, w)
}
