package transcriber

import (
	"github.com/xlab/pocketsphinx-go/sphinx"
	"os"
	"testing"
)

func TestTranscribe(t *testing.T) {
	want := "this is a test"

	cfg := sphinx.NewConfig(
		sphinx.HMMDirOption("/usr/local/share/pocketsphinx/model/en-us/en-us"),
		sphinx.DictFileOption("/usr/local/share/pocketsphinx/model/en-us/cmudict-en-us.dict"),
		sphinx.LMFileOption("/usr/local/share/pocketsphinx/model/en-us/en-us.lm.bin"),
		sphinx.SampleRateOption(16000),
		sphinx.LogFileOption("sphinx.log"),
	)

	dec, err := sphinx.NewDecoder(cfg)
	if err != nil {
		t.Error(err)
	}

	file, err := os.Open("audio/test3_en.wav")
	//file, err := os.Open("audio/test4_2_en.wav")
	if err != nil {
		t.Error(err)
	}
	defer file.Close()

	c := ReadFile(file)
	got, _, err := Transcribe(dec, c)
	if err != nil {
		t.Error(err)
	}
	if got != want {
		t.Errorf("Transcribe(dec, c) = %q, want %q", got, want)
	}
}

