package handler

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/vbelorus/go-app/v2/src/config"
	"github.com/vbelorus/go-app/v2/src/models"
	"net"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type DeviceEventHandler struct {
	App                *config.Application
	DeviceEventChannel chan models.DeviceEvent
}

func (h *DeviceEventHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Add("Allow", http.MethodPost)
		h.App.ClientError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	decoder := json.NewDecoder(r.Body)
	validate := validator.New()
	ip, errIp := getIP(r)

	withErrors := false
	savedSuccess := 0
	for decoder.More() {
		var event models.DeviceEvent
		if err := decoder.Decode(&event); err != nil {
			withErrors = true
			h.App.ClientError(w, http.StatusBadRequest, fmt.Sprintf("JSON decode error: %v", err))
			//return if it's *json.SyntaxError, continue for another errors
			switch reflect.TypeOf(err).String() {
			case "*json.SyntaxError":
				fmt.Fprintln(w, fmt.Sprintf("Successfully saved: %d", savedSuccess))
				return
			default:
				continue
			}
		}

		event.Ip = ip
		if errIp != nil {
			event.Ip = "No valid ip"
		}
		event.ServerTime = time.Now()

		err := validate.Struct(event)
		if err != nil {
			withErrors = true
			h.App.ClientError(w, http.StatusBadRequest, fmt.Sprintf("Entity validate error: %v", err))
			continue
		}

		h.App.DeviceEventChannel <- event
		savedSuccess++
	}

	if withErrors {
		fmt.Fprintln(w, fmt.Sprintf("Successfully saved: %d", savedSuccess))
	}
}

func getIP(r *http.Request) (string, error) {
	//Get IP from the X-REAL-IP header
	ip := r.Header.Get("X-REAL-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}

	//Get IP from X-FORWARDED-FOR header
	ips := r.Header.Get("X-FORWARDED-FOR")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP := net.ParseIP(ip)
		if netIP != nil {
			return ip, nil
		}
	}

	//Get IP from RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	netIP = net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}
	return "", fmt.Errorf("No valid ip found")
}
