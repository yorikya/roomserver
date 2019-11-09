package httpserver

import (
	"encoding/json"
	"net/http"

	"github.com/yorikya/roomserver/client"
	"github.com/yorikya/roomserver/mqttserver"
)

func InitRoutes(s *mqttserver.Server) {
	http.HandleFunc("/data", withServerData(s))

	// This works and strip "/static/" fragment from path
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

}

func withServerData(s *mqttserver.Server) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		allClients := []client.Client{}
		for _, v := range s.GetClients() {
			allClients = append(allClients, *v)
		}
		js, err := json.Marshal(allClients)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}
