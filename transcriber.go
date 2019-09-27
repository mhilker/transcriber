package transcriber

import (
	"encoding/binary"
	"errors"
	"github.com/xlab/pocketsphinx-go/sphinx"
	"io"
	"os"
)

func Transcribe(dec *sphinx.Decoder, c <-chan []int16) (string, int32, error) {
	if !dec.StartUtt() {
		return "", 0, errors.New("sphinx failed to start utterance")
	}
	for data := range c {
		dec.ProcessRaw(data, false, false)
	}
	if !dec.EndUtt() {
		return "", 0, errors.New("sphinx failed to stop utterance")
	}
	hyp, score := dec.Hypothesis()
	return hyp, score, nil
}

func ReadFile(file *os.File) <-chan []int16 {
	c := make(chan []int16)
	go func() {
		data := make([]int16, 512)
		for {
			err := binary.Read(file, binary.LittleEndian, data)
			if err == io.EOF {
				break
			}
			if err == io.ErrUnexpectedEOF {
				c <- data
				break
			}
			c <- data
		}
		close(c)
	}()
	return c
}