package main

// typedef unsigned char Uint8;
// void SineWave(void *userdata, Uint8 *stream, int len);
import "C"
import (
	"io/ioutil"
	"log"
	"math"
	"reflect"
	"unsafe"

	"github.com/veandco/go-sdl2/mix"
)

const (
	toneHz   = 440
	sampleHz = 2400
	dPhase   = 2 * math.Pi * toneHz / sampleHz
)

//export SineWave
func SineWave(userdata unsafe.Pointer, stream *C.Uint8, length C.int) {
	n := int(length)
	hdr := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(stream)), Len: n, Cap: n}
	buf := *(*[]C.Uint8)(unsafe.Pointer(&hdr))

	var phase float64
	for i := 0; i < n; i += 2 {
		phase += dPhase
		sample := C.Uint8((math.Sin(phase) + 0.999999) * 128)
		buf[i] = sample
		buf[i+1] = sample
	}
}

// RunAudioHandler takes a bool channel, which turns audio on and off
func RunAudioHandler(audio <-chan bool) {

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
