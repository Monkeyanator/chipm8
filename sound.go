package main

// typedef unsigned char Uint8;
// void SineWave(void *userdata, Uint8 *stream, int len);
import "C"
import (
	"io/ioutil"
	"log"

	"github.com/veandco/go-sdl2/mix"
)

// RunAudioHandler takes a bool channel, which turns audio on and off
func (chip *chip8) RunAudioHandler(audio <-chan bool) {

	if err := mix.OpenAudio(44100, mix.DEFAULT_FORMAT, 2, 4096); err != nil {
		log.Println(err)
		return
	}

	// Load BEEP wav
	data, err := ioutil.ReadFile("./assets/beep.wav")
	if err != nil {
		log.Println(err)
	}

	defer mix.CloseAudio()
	activated := false

	for {
		shouldActivate := <-audio
		if activated == shouldActivate {
			continue
		}

		if shouldActivate {

			// Load sound again, and go for it
			chunk, err := mix.QuickLoadWAV(data)
			if err != nil {
				log.Println(err)
			}
			defer chunk.Free()

			chunk.Play(1, 1)
		}

		activated = shouldActivate
	}

}
