package api

import (
	"context"
	"encoding/json"
	"iotstarter/internal/auth"
	"iotstarter/internal/broker"
	"iotstarter/internal/config"
	"iotstarter/internal/model"
	"iotstarter/internal/store"
	"iotstarter/internal/views"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	store  store.Store
	broker broker.Broker
	views  *views.Views
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) WithStore(store store.Store) *Handler {
	h.store = store
	return h
}

func (h *Handler) WithBroker(broker broker.Broker) *Handler {
	h.broker = broker
	return h
}

func (h *Handler) WithViews(views *views.Views) *Handler {
	h.views = views
	return h
}

func (h *Handler) registerUserRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	if h.store != nil {
		mux.HandleFunc("POST /api/auth/register", h.registerUser)
		mux.HandleFunc("POST /api/auth/login", h.logInUser)
		mux.HandleFunc("POST /api/auth/logout", h.logOutUser)
		mux.HandleFunc("POST /api/devices", h.registerDevice)
		mux.HandleFunc("GET /api/devices", h.getDevices)
		mux.HandleFunc("PATCH /api/devices/{id}/reauth", h.reauthDevice)
		mux.HandleFunc("GET /api/devices/{id}", h.getDeviceById)
		mux.HandleFunc("GET /api/devices/{id}/measurements", h.getDeviceMeasurements)
	}

	if h.broker != nil {
		mux.HandleFunc("POST /api/measurements", h.saveMeasurement)
	}

	if h.views != nil {
		mux.Handle("GET /static/", http.StripPrefix("/static/", http.HandlerFunc(h.loadStaticResources)))
		mux.HandleFunc("GET /", h.getHomePage)
		mux.HandleFunc("GET /devices", h.getDevicesPage)
	}

	return mux
}

func (h *Handler) getHomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	// http.ServeFile(w, r, "static/html/index.html")
	err := h.views.RenderPage(w, r, "home", nil)
	if err != nil {
		log.Println("error getting home view", err.Error())
		return
	}
}

func (h *Handler) getDevicesPage(w http.ResponseWriter, r *http.Request) {
	type deviceT struct {
		Name string
		Type string
	}
	type devicesT struct {
		Devices []deviceT
	}

	devices := devicesT{
		Devices: []deviceT{{Name: "device1", Type: "sensor"}, {Name: "device2", Type: "light"}},
	}
	err := h.views.RenderPage(w, r, "devices", devices)
	if err != nil {
		log.Println("error getting rendering page", err.Error())
		return
	}
}

func (h *Handler) loadStaticResources(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "public, max-age=31536000")
	fs := http.FileServer(http.Dir("static"))
	fs.ServeHTTP(w, r)
}

func (h *Handler) registerUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*3))
	defer cancel()

	username := r.FormValue("username")
	password := r.FormValue("password")
	user, err := h.store.RegisterUser(ctx, username, password)

	if err == store.ErrUsernameTaken {
		http.Error(w, "username taken", http.StatusConflict)
		return
	}

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "error registering user", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) logInUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*3))
	defer cancel()

	username := r.FormValue("username")
	password := r.FormValue("password")
	user, err := h.store.GetUserFromCreds(ctx, username, password)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "invalid credentials", http.StatusForbidden)
		return
	}

	sesh, err := h.store.CreateUserSession(ctx, user.ID)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "error logging in", http.StatusInternalServerError)
		return
	}

	auth.SetCookie(w, sesh.Token)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) logOutUser(w http.ResponseWriter, r *http.Request) {
	auth.ClearCookie(w)
	http.Redirect(w, r, "/", http.StatusAccepted)
}

func (h *Handler) getDeviceMeasurements(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*3))
	defer cancel()

	params, err := getMeasurementsQueryParams(r)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	measurements, err := h.store.GetDeviceMeasurements(ctx, params.deviceId, params.start, params.end)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Error retrieving measurements", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(measurements)
}

func (h *Handler) registerDevice(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*3))
	defer cancel()

	location := r.FormValue("location")
	if location == "" {
		http.Error(w, "location required", http.StatusBadRequest)
		return
	}

	device, err := h.store.RegisterDevice(ctx, 1, location)
	if err != nil {
		log.Println("ERROR:", err.Error())
		http.Error(w, "Error registering device", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(device)
}

func (h *Handler) reauthDevice(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*3))
	defer cancel()

	deviceId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid device Id", http.StatusBadRequest)
		return
	}
	device, err := h.store.ReauthDevice(ctx, 1, deviceId)
	if err != nil {
		log.Println("ERROR:", err.Error())
		http.Error(w, "Error reauthing device", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(device)
}

func (h *Handler) getDevices(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*3))
	defer cancel()

	devices, err := h.store.GetDevices(ctx)
	if err != nil {
		log.Println("ERROR:", err.Error())
		http.Error(w, "Error retrieving devices", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(devices)
}

func (h *Handler) getDeviceById(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*3))
	defer cancel()

	deviceId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid device Id", http.StatusBadRequest)
		return
	}

	device, err := h.store.GetDeviceById(ctx, deviceId)
	if err == store.ErrDeviceNotFound {
		http.Error(w, "Device not found", http.StatusNotFound)
		return
	}

	if err != nil {
		log.Println("ERROR:", err.Error())
		http.Error(w, "Error retrieving devices", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(device)
}

func (h *Handler) saveMeasurement(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	measurement := &model.Measurement{}
	if err := json.NewDecoder(r.Body).Decode(&measurement); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.broker.Publish(config.BrokerMeasurementSubject, measurement); err != nil {
		http.Error(w, "Failed to publish", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
