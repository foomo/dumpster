package handler

import (
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
}

// Restore a dump
func (restore *Restore) Restore(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	restoreReport, restoreErrors, err := restore.dumpster.RestoreDump(ps.ByName("type"), ps.ByName("id"))
	if err != nil {
		errReply(w, http.StatusInternalServerError, err)
		return
	}
	jsonReply(&responses.RestoreReport{
		Report: restoreReport,
		Errors: restoreErrors,
	}, w)

}
