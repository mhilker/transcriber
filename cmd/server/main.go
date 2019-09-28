package main

import (
	"bytes"
	"encoding/json"
	"github.com/mhilker/transcriber"
	"github.com/xlab/pocketsphinx-go/sphinx"
	"log"
	"net/http"
	"os"
	"os/exec"
)

type successResponse struct {
	Hypothesis string  `json:"hypothesis"`
	Score      float64 `json:"score"`
}

type errorResponse struct {
	Message string `json:"message"`
}

func main() {
	log.Println("loading sphinx decoder")
	dec, err := newDecoder()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("handling request")
		httpHandler(w, r, dec)
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
		sphinx.InputEndianOption("little"),
		sphinx.LogFileOption("sphinx.log"),
	)

	return sphinx.NewDecoder(cfg)
}

func httpHandler(w http.ResponseWriter, r *http.Request, dec *sphinx.Decoder) {
	if r.Method != http.MethodPost {
		httpError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	if r.Header.Get("Content-Type") != "audio/webm" {
		httpError(w, http.StatusInternalServerError, "invalid content type")
		return
	}

	buf := bytes.Buffer{}
	cmd := exec.Command("/usr/bin/ffmpeg", "-i", "-", "-acodec", "pcm_s16le", "-ar", "16000", "-f", "s16le", "-")
	cmd.Stdin = r.Body
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		httpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	c := transcriber.Read(&buf)
	hyp, score, err := transcriber.Transcribe(dec, c)
	if err != nil {
		httpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := successResponse{hyp, score}
	j, err := json.Marshal(resp)
	if err != nil {
		httpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(j)
}

func httpError(w http.ResponseWriter, status int, message string) {
	resp := errorResponse{message}
	j, _ := json.Marshal(resp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(j)
}
