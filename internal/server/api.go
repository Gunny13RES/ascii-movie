package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

type ApiServer struct {
	Server
	TelnetEnabled bool
	SSHEnabled    bool
}

func NewApi(flags *flag.FlagSet) ApiServer {
	return ApiServer{Server: NewServer(flags, ApiFlagPrefix)}
}

func (s *ApiServer) Listen(ctx context.Context) error {
	s.Log.WithField("address", s.Address).Info("Starting API server")

	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.Status)
	server := http.Server{Addr: s.Address, Handler: mux}
	go func() {
		<-ctx.Done()
		s.Log.Info("Stopping API server")
		defer s.Log.Info("Stopped API server")
		if err := server.Close(); err != nil {
			log.WithError(err).Error("Failed to close server")
		}
	}()

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

type StatusResponse struct {
	Healthy bool `json:"healthy"`
	SSH     bool `json:"ssh"`
	Telnet  bool `json:"telnet"`
}

func (s *ApiServer) Status(w http.ResponseWriter, r *http.Request) {
	response := StatusResponse{
		Telnet: telnetListeners == 1,
		SSH:    sshListeners == 1,
	}
	response.Healthy = (!s.SSHEnabled || response.SSH) && (!s.TelnetEnabled || response.Telnet)

	buf, err := json.Marshal(response)
	if err != nil {
		s.Log.WithError(err).Error("Failed to marshal API response")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if response.Healthy {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	if _, err := w.Write([]byte(buf)); err != nil {
		s.Log.WithError(err).Error("Failed to write API response")
	}
}
