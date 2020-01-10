package httphandlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/yorikya/roomserver/client"
	"github.com/yorikya/roomserver/devices"
	"github.com/yorikya/roomserver/server"

	"github.com/gorilla/websocket"
)

//HTML Tmplates directory
var (
	templates = template.Must(template.ParseGlob("templates/*"))

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
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
			if device := c.GetDeviceByName(deviceID); device != nil {
				action, err := getURLParam(r, "action")
				if err != nil {
					log.Println(err)
					fmt.Fprintln(w, err)
					return
				}
				cmd, val, err := device.CreateCMD(action)
				if err != nil {
					log.Println(err)
					fmt.Fprintln(w, err)
					return
				}

				url := fmt.Sprintf("http://%s/action?deviceid=%s&val=%s&cmd=%s", c.IPstr, deviceID, val, cmd)
				log.Println("the client action url:", url)
				res, err := http.Get(url)
				if err != nil {
					log.Printf("get an erro when send command cto client: %s, cmd: %s\n", roomID, cmd)
					fmt.Fprintln(w, err)
					return
				}

				device.SetValue(val)
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

func withServerSelectRoom(s *server.Server) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		name, err := getURLParam(r, "name")
		if err != nil {
			log.Println(err)
			fmt.Fprintln(w, err)
			return
		}
		mainRomm := fmt.Sprintf("%s_main", name)
		if c := s.GetClient(mainRomm); c != nil {
			cameraID := devices.Camera2MP
			camroom := fmt.Sprintf("%s_cam", name)
			clientCam := s.GetClient(camroom)
			if clientCam == nil {
				log.Printf("room '%s' does not have cammera (%s), init default camera", name, camroom)
				clientCam = client.NewClient(camroom, devices.NewCamera(camroom, cameraID))
			}
			rgbStripID := devices.RGBstrip
			irac := devices.IR_ac_aircool
			dhth := devices.DHT_Humidity
			dhtt := devices.DHT_Temperature
			ac := strings.Split(c.GetDeviceByName(irac).GetValueStr(), ",")
			d := struct {
				RoomID,
				DHTHumuditi,
				DHTTemperture,
				CameraID,
				DHTSensorHumudutyHTML,
				DHTSensorTempertureHTML,
				RGBStripVal,
				ACModeVal,
				ACTempertureVal,
				CameraStatus,
				ACName,
				RGBName string

				RGBOptions,
				ACModeOptions,
				ACTempertureOptions []string
			}{ //For device name see rooms.json config
				RGBName:                 rgbStripID,
				ACName:                  irac,
				DHTHumuditi:             dhth,
				DHTTemperture:           dhtt,
				CameraID:                cameraID,
				RoomID:                  mainRomm,
				DHTSensorHumudutyHTML:   fmt.Sprintf("<h2 class='text-success'>%s</h2>", c.GetDeviceByName(dhth).GetValueStr()),
				DHTSensorTempertureHTML: fmt.Sprintf("<h2 class='text-success'>%s</h2>", c.GetDeviceByName(dhtt).GetValueStr()),
				RGBStripVal:             c.GetDeviceByName(rgbStripID).GetValueStr(),
				ACModeVal:               ac[0],
				ACTempertureVal:         ac[1],
				ACModeOptions:           c.GetDeviceByName(irac).GetOptions("mode"),
				ACTempertureOptions:     c.GetDeviceByName(irac).GetOptions("temp"),
				RGBOptions:              c.GetDeviceByName(rgbStripID).GetOptions(""),
				CameraStatus:            fmt.Sprintf("<h2 class='text-success'>%s</h2>", clientCam.GetDeviceByName(cameraID).GetValueStr()),
			}
			err := templates.ExecuteTemplate(w, "room.html", d) //execute the template and pass it the HomePageVars struct to fill in the gaps
			if err != nil {                                     // if there is an error
				log.Print("template executing error: ", err) //log it
				fmt.Fprintln(w, err)
			}
			log.Println("generate room.html, room name", name)
		} else {
			log.Println("client", name, "not exists")
			fmt.Fprintln(w, "can not find client ", name)
		}
	}
}

func withServerRooms(s *server.Server) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		err := templates.ExecuteTemplate(w, "rooms.html", s.GetRooms()) //execute the template and pass it the HomePageVars struct to fill in the gaps
		if err != nil {                                                 // if there is an error
			log.Print("template executing error: ", err) //log it
		}

	}
}

func withServerWS(s *server.Server) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//TODO: Store connection
		log.Printf("get ws connection: %+v\n", r)

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		s.RoomHub.Register(conn)
	}
}

func serverLogin(w http.ResponseWriter, r *http.Request) {
	log.Printf("get login request: %+v", r)
	if r.Method == "GET" {

		err := templates.ExecuteTemplate(w, "login.html", nil) //execute the template and pass it the HomePageVars struct to fill in the gaps
		if err != nil {                                        // if there is an error
			log.Print("template executing error: ", err) //log it
		}

	} else {
		r.ParseForm()
		// logic part of log in
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
		http.Redirect(w, r, "/rooms", http.StatusSeeOther)
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
			value, err := getURLParam(r, "value")
			if err != nil {
				log.Println(err)
				fmt.Fprintln(w, err)
				return
			}

			c.Update(device, value)
			if err = s.RoomHub.Brodcast(fmt.Sprintf("%s/update/%s/text-success/%s", clientid, device, value)); err != nil {
				log.Println("failed broadcast message")
				fmt.Fprintln(w, err)
				return
			}
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
			Success bool
			ErrMsg  string
		}

		if c := s.GetClient(clientid); c != nil {
			c.UpdateIPstr(cutClientIP(r.RemoteAddr))
			resp := AuthResponse{
				Success: true,
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
			Success: false,
			ErrMsg:  msg,
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
	http.HandleFunc("/update", withServerUpdate(s))
	http.HandleFunc("/auth", withServerAuth(s))

	//User URL
	http.HandleFunc("/data", withServerData(s))
	http.HandleFunc("/action", withServerAction(s))

	http.HandleFunc("/login", serverLogin)
	http.HandleFunc("/rooms", withServerRooms(s))
	http.HandleFunc("/room", withServerSelectRoom(s))
	http.HandleFunc("/ws", withServerWS(s))

	// This works and strip "/static/" fragment from path
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

}
