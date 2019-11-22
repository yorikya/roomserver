package httpserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/yorikya/roomserver/client"
	"github.com/yorikya/roomserver/mqttserver"
)

func getURLParam(r *http.Request, key string) (string, error) {
	keys, ok := r.URL.Query()[key]
	if !ok || len(keys[0]) < 1 {
		return "", fmt.Errorf("Url Param '%s' is missing", key)
	}
	return keys[0], nil
}

func withServerAction(s *mqttserver.Server) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// url example http://localhost:3000/action?roomid=office&deviceid=AC&action=on
		roomID, err := getURLParam(r, "roomid")
		if err != nil {
			log.Println(err)
			fmt.Fprintln(w, err)
			return
		}
		deviceID, err := getURLParam(r, "deviceid")
		if err != nil {
			log.Println(err)
			fmt.Fprintln(w, err)
			return
		}
		action, err := getURLParam(r, "action")

		s.Publish(fmt.Sprintf("%sInTopic", roomID), []byte(fmt.Sprintf("/%s/%s", deviceID, action)))
		log.Printf("get an request, room-id: %s, device-id: %s, action: %s", roomID, deviceID, action)
	}
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

func InitRoutes(s *mqttserver.Server) {
	http.HandleFunc("/data", withServerData(s))
	http.HandleFunc("/action", withServerAction(s))

	// This works and strip "/static/" fragment from path
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

}
