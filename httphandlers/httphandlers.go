package httphandlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/yorikya/roomserver/client"
	"github.com/yorikya/roomserver/devices"
	"github.com/yorikya/roomserver/server"
)

func getURLParam(r *http.Request, key string) (string, error) {
	keys, ok := r.URL.Query()[key]
	if !ok || len(keys[0]) < 1 {
		return "", fmt.Errorf("Url Param '%s' is missing", key)
	}
	return keys[0], nil
}

func withServerAction(s *server.Server) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// url example http://localhost:3000/action?roomid=room1_main&deviceid=rgbstrip&action=1900L
		log.Println("get an /action:", r)
		roomID, err := getURLParam(r, "roomid")
		if err != nil {
			log.Println(err)
			fmt.Fprintln(w, err)
			return
		}
		if c := s.GetClient(roomID); c != nil {
			deviceID, err := getURLParam(r, "deviceid")
			if err != nil {
				log.Println(err)
				fmt.Fprintln(w, err)
				return
			}
			if s := c.GetSensorByName(deviceID); s != nil {
				action, err := getURLParam(r, "action")
				if err != nil {
					log.Println(err)
					fmt.Fprintln(w, err)
					return
				}
				cmd, err := s.CreateCMD(action)
				if err != nil {
					log.Println(err)
					fmt.Fprintln(w, err)
					return
				}
				url := fmt.Sprintf("http://%s/action?deviceid=%s&cmd=%s", c.IPstr, deviceID, cmd)
				log.Println("the client action url:", url)
				res, err := http.Get(url)
				if err != nil {
					log.Printf("get an erro when send command cto client: %s, cmd: %s\n", roomID, cmd)
				}

				fmt.Fprintln(w, err)
				log.Println("the response from client", res)
				return
			}
			log.Println("clientID: %s does not have device: %s\n", roomID, deviceID)
			return
		}

		log.Printf("/action does not match clientID: %s\n", roomID)
	}
}

func withServerData(s *server.Server) func(w http.ResponseWriter, r *http.Request) {
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

func cutClientIP(ip string) string {
	return strings.SplitN(ip, ":", 2)[0]
}

// /update?device=hdt&sensor=Temperature&value=23.70&clientid=room1_hdt
func withServerUpdate(s *server.Server) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("get /update request: %+v\n", r)
		clientid, err := getURLParam(r, "clientid")
		if err != nil {
			log.Println(err)
			fmt.Fprintln(w, err)
			return
		}

		if c := s.GetClient(clientid); c != nil {
			c.UpdateIPstr(cutClientIP(r.RemoteAddr))
			device, err := getURLParam(r, "device")
			if err != nil {
				log.Println(err)
				fmt.Fprintln(w, err)
				return
			}
			sensor, err := getURLParam(r, "sensor")
			if err != nil {
				log.Println(err)
				fmt.Fprintln(w, err)
				return
			}
			value, err := getURLParam(r, "value")
			if err != nil {
				log.Println(err)
				fmt.Fprintln(w, err)
				return
			}

			c.Update(device, sensor, value)
			return
		}
		err = fmt.Errorf("server dos not have clientID: %s", clientid)
		log.Println(err)
		fmt.Fprintln(w, err)
		return
	}
}

func withServerAuth(s *server.Server) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("get /auth request: %+v\n", r)
		clientid, err := getURLParam(r, "clientid")
		if err != nil {
			log.Println(err)
			fmt.Fprintln(w, err)
			return
		}

		type AuthResponse struct {
			Success    bool
			Devices    map[string]devices.Device
			DevicesNum int
			ErrMsg     string
		}
		m := make(map[string]devices.Device)
		if c := s.GetClient(clientid); c != nil {
			c.UpdateIPstr(cutClientIP(r.RemoteAddr))
			for _, d := range c.Devices {
				m[fmt.Sprintf("%s_%s", d.GetID(), d.GetSensor())] = d
			}
			resp := AuthResponse{
				Success:    true,
				Devices:    m,
				DevicesNum: len(m),
			}
			b, err := json.Marshal(resp)
			if err != nil {
				log.Println("error:", err)
			}
			w.Write(b)
			log.Printf("reponse auht: %s", string(b))
			return
		}
		msg := fmt.Sprintf("authetication failed no have clientID: %s", clientid)
		resp := AuthResponse{
			Success:    false,
			Devices:    m,
			DevicesNum: len(m),
			ErrMsg:     msg,
		}
		b, err := json.Marshal(resp)
		if err != nil {
			log.Println("error:", err)
		}
		w.Write(b)
		log.Println(msg)
		return
	}
}

func InitRoutes(s *server.Server) {
	http.HandleFunc("/data", withServerData(s))
	http.HandleFunc("/action", withServerAction(s))
	http.HandleFunc("/update", withServerUpdate(s))
	http.HandleFunc("/auth", withServerAuth(s))

	// This works and strip "/static/" fragment from path
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

}
