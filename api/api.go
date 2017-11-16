package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/tidwall/gjson"
)

const (
	CT        = "Content-Type"
	MIME_JSON = "application/json"
)

type InputData struct {
	Ifa      string `json:"ifa"`
	Country  string `json:"country"`
	App      string `json:"app"`
	Platform string `json:"platform"`
}

type Config struct {
	ServicePort       string        `json:"service_port"`
	RedisHost         string        `json:"redis_host"`
	RedisPort         string        `json:"redis_port"`
	MaxIddleCons      int           `json:"max_iddle_connections"`
	IFACounterTTL     int           `json:"ifa_counter_ttl"`
	IFASeriesInterval time.Duration `json:"ifa_series_interval"`
}

// stat структура для описания одной строки старистики
type stat struct {
	Country  string `json:"country"`
	App      string `json:"app"`
	Platform string `json:"platform"`
	Count    int    `json:"count"`
}

type Api struct {
	*log.Logger
	*service
	port string
}

func CreateApi(config *Config, logger *log.Logger) *Api {
	api := Api{
		Logger:  logger,
		service: CreateService(config),
		port:    config.ServicePort,
	}

	return &api
}

func (api *Api) Run() {
	http.HandleFunc("/", api.inputHandler)
	http.HandleFunc("/stats", api.statsHandler)
	api.Fatal(http.ListenAndServe(":"+api.port, nil))
}

func (api *Api) inputHandler(w http.ResponseWriter, r *http.Request) {

	if ct := r.Header.Get("Content-Type"); ct != MIME_JSON {
		http.Error(w, "wrong content type", http.StatusBadRequest)
		return
	}

	bts, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		http.Error(w, "wrong body", http.StatusBadRequest)
		return
	}

	ifa, stat := parseInput(string(bts))
	pos, err := api.Process(ifa, stat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set(CT, MIME_JSON)
	fmt.Fprintf(w, `{"pos": %d}`, pos)
}

func (api *Api) statsHandler(w http.ResponseWriter, r *http.Request) {

	stats, err := api.Stats()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set(CT, MIME_JSON)
	json.NewEncoder(w).Encode(stats)
}

// parseInput получение необходимых значений из входного json
func parseInput(json string) (ifa, statKey string) {
	res := gjson.GetMany(json, "device.ifa", "geo.country", "app.name", "device.os")
	ifa = res[0].String()

	buf := bytes.NewBufferString("{")
	buf.WriteString(`"country":"` + res[1].String() + `",`)
	buf.WriteString(`"app":"` + res[2].String() + `",`)
	buf.WriteString(`"platform":"` + res[1].String() + `"}`)
	statKey = buf.String()

	return
}
