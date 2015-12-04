package server

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/abbot/go-http-auth"
	"github.com/foomo/dumpster/config"
	"github.com/foomo/dumpster/dumpster"
	"github.com/foomo/dumpster/server/handler"
	"github.com/julienschmidt/httprouter"
)

// Server dumpster server
type Server struct {
	Config   *config.Config
	Router   *httprouter.Router
	Dumpster *dumpster.Dumpster
}

func newServer(c *config.Config) (s *Server, err error) {
	d, err := dumpster.NewDumpster(c.DataDir, c.Dumps, c.Remotes)
	if err != nil {
		return nil, err
	}
	s = &Server{
		Config:   c,
		Router:   httprouter.New(),
		Dumpster: d,
	}
	handler.SetupHandlers(s.Router, s.Dumpster)
	return s, nil
}

type BasicAuthHandler struct {
	Server        *Server
	authenticator *auth.BasicAuth
}

func NewBasicAuthHandler(server *Server, htpasswordFile string) (ba *BasicAuthHandler) {
	secretProvider := auth.HtpasswdFileProvider(htpasswordFile)
	authenticator := auth.NewBasicAuthenticator("dumpster", secretProvider)
	return &BasicAuthHandler{
		Server:        server,
		authenticator: authenticator,
	}
}

func (ba *BasicAuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user := ba.authenticator.CheckAuth(r)
	if len(user) == 0 {
		ba.authenticator.RequireAuth(w, r)
		return
	}
	ba.Server.Router.ServeHTTP(w, r)
}

func getTLSConfig() *tls.Config {
	c := &tls.Config{}
	c.MinVersion = tls.VersionTLS12
	c.PreferServerCipherSuites = true
	c.CipherSuites = []uint16{
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
	}
	c.CurvePreferences = []tls.CurveID{
		tls.CurveP256,
		tls.CurveP384,
		tls.CurveP521,
	}
	return c
}

// Run run as a server
func Run(c *config.Config) error {
	s, err := newServer(c)
	if err != nil {
		return err
	}
	log.Println("Starting Dumpster Server on: ", c.HTTP.Address)
	log.Println("  basic auth from: ", c.HTTP.BasicAuthFile)
	log.Println("  data dir:", c.DataDir)
	ba := NewBasicAuthHandler(s, c.HTTP.BasicAuthFile)
	errorChan := make(chan (error))
	if len(c.HTTP.Address) > 0 {
		go func() {
			errorChan <- http.ListenAndServe(c.HTTP.Address, ba)
		}()
	}
	if c.HTTP.TLS != nil {
		go func() {
			log.Println("tls is configured: ", c.HTTP.TLS)
			tlsServer := &http.Server{
				Addr:      c.HTTP.TLS.Address,
				Handler:   ba,
				TLSConfig: getTLSConfig(),
			}
			errorChan <- tlsServer.ListenAndServeTLS(c.HTTP.TLS.Cert, c.HTTP.TLS.Key)
		}()
	}
	return <-errorChan
}
