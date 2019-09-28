package main

import (
	"encoding/json"
	"github.com/mhilker/transcriber"
	"github.com/xlab/pocketsphinx-go/sphinx"
	"log"
	"net/http"
)

func main() {
	log.Println("loading sphinx decoder")
	dec, err := newDecoder()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("handling request")
		handler(w, r, dec)
	})

	log.Println("starting server")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func newDecoder() (*sphinx.Decoder, error) {
	cfg := sphinx.NewConfig(
		sphinx.HMMDirOption("/usr/local/share/pocketsphinx/model/en-us/en-us"),
		sphinx.DictFileOption("/usr/local/share/pocketsphinx/model/en-us/cmudict-en-us.dict"),
		sphinx.LMFileOption("/usr/local/share/pocketsphinx/model/en-us/en-us.lm.bin"),
		sphinx.SampleRateOption(16000),
		sphinx.LogFileOption("sphinx.log"),
	)

	return sphinx.NewDecoder(cfg)
}

type successResponse struct {
	Hypothesis string `json:"hypothesis"`
	Score      int32  `json:"score"`
}

type errorResponse struct {
	Message string `json:"message"`
}

func handler(w http.ResponseWriter, r *http.Request, dec *sphinx.Decoder) {
	if r.Method != http.MethodPost {
		resp := errorResponse{"method not allowed"}
		j, _ := json.Marshal(resp)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(j)
		return
	}

	c := transcriber.Read(r.Body)
	hyp, score, err := transcriber.Transcribe(dec, c)
	if err != nil {
		resp := errorResponse{err.Error()}
		j, _ := json.Marshal(resp)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(j)
		return
	}

	resp := successResponse{hyp, score}
	j, err := json.Marshal(resp)
	if err != nil {
		resp := errorResponse{err.Error()}
		j, _ := json.Marshal(resp)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(j)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(j)
}
