package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/veandco/go-sdl2/mix"
)

var chunkLaser *mix.Chunk
var chunkHit *mix.Chunk
var chunkExplosion *mix.Chunk
var music1 *mix.Music
var music2 *mix.Music

func initSounds() {

	fmt.Println("Init sounds")

	if err := mix.OpenAudio(44100, mix.DEFAULT_FORMAT, 4, 4096); err != nil {
		log.Println(err)
		return
	}
	//defer mix.CloseAudio()

	// Load entire WAV data from file
	dataLaser, err := ioutil.ReadFile("sounds/laser.wav")
	if err != nil {
		log.Println(err)
	}

	// Load WAV from data (memory)
	chunkLaser, err = mix.QuickLoadWAV(dataLaser)
	if err != nil {
		log.Println(err)
	}
	//defer chunk.Free()

	// Play 4 times
	//chunk.Play(1, 3)

	// Wait until it finishes playing
	//for mix.Playing(-1) == 1 {
	//	sdl.Delay(16)
	//}

	dataHit, err := ioutil.ReadFile("sounds/hit.wav")
	if err != nil {
		log.Println(err)
	}

	// Load WAV from data (memory)
	chunkHit, err = mix.QuickLoadWAV(dataHit)
	if err != nil {
		log.Println(err)
	}

	dataExplosion, err := ioutil.ReadFile("sounds/explosion.wav")
	if err != nil {
		log.Println(err)
	}

	// Load WAV from data (memory)
	chunkExplosion, err = mix.QuickLoadWAV(dataExplosion)
	if err != nil {
		log.Println(err)
	}

	//music, err = mix.LoadMUS("sounds/test.mp3") // or music.wav
	music1, err = mix.LoadMUS("sounds/bestmid3/127-yell.mid")
	if err != nil {
		log.Println(err)
	}
	music2, err = mix.LoadMUS("sounds/test.mp3")
	if err != nil {
		log.Println(err)
	}
	//err = music1.Play(-1)
	if err != nil {
		log.Println(err)
	}

	//	sdl.Delay(10000)

	//	err = music1.Play(3)
	//	if err != nil {
	//		log.Println(err)
	//	}
}
