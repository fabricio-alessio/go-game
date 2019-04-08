package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/veandco/go-sdl2/mix"
)

const (
	channelsQuantity = 16
	maxPlayCount     = 100
)

var music1 *mix.Music
var music2 *mix.Music

func initSounds() {

	fmt.Println("Init sounds")

	if err := mix.OpenAudio(mix.DEFAULT_FREQUENCY, mix.DEFAULT_FORMAT, 2, 4096); err != nil {
		log.Println(err)
		return
	}
	//defer mix.CloseAudio()

	mix.AllocateChannels(channelsQuantity)

	//	chunkLaser = loadSound2("sounds/laser.wav")
	//	chunkHit = loadSound2("sounds/hit.wav")
	//	chunkExplosion = loadSound2("sounds/explosion.wav")
	//defer chunk.Free()

	var err error
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

func loadSound(fileName string) ([]byte, *mix.Chunk) {

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Println(err)
	}

	chunk := loadSoundFromData(data)

	return data, chunk
}

func loadSoundFromData(data []byte) *mix.Chunk {

	chunk, err := mix.QuickLoadWAV(data)
	fmt.Printf("Sound loaded2 %v\n", chunk)
	if err != nil {
		log.Println(err)
	}
	return chunk
}

func playSound(chunk *mix.Chunk) int {

	played, err := chunk.Play(-1, 0)
	//fmt.Printf("| channel %2v | hit %8d | laser %8d | explosion %8d | all %8d |\n",
	//	played, countHit, countLaser, countExplosion, countHit+countLaser+countExplosion)
	if err != nil {
		fmt.Printf("Error playing sound %s\n", err)
	}
	return played
}

type sound struct {
	countPlays  int
	soundData   []byte
	chunk       *mix.Chunk
	fileName    string
	lastChannel int
}

func newSound(fileName string) *sound {

	s := sound{
		countPlays: 0,
		fileName:   fileName}

	s.soundData, s.chunk = loadSound(fileName)
	return &s
}

func (s *sound) play() {

	//fmt.Printf("Will play sound %s count %d\n", s.fileName, s.countPlays)
	if s.countPlays > maxPlayCount {
		s.reset()
	} else {
		s.countPlays++
		s.lastChannel = playSound(s.chunk)
	}
}

func (s *sound) reset() {

	isPlaying := mix.Playing(s.lastChannel)
	aChunk := mix.GetChunk(s.lastChannel)
	if isPlaying == 1 && aChunk == s.chunk {
		fmt.Printf("Playing sound %s, stopping before reset\n", s.fileName)
		mix.HaltChannel(s.lastChannel)
	} else {
		s.chunk.Free()
		s.chunk = loadSoundFromData(s.soundData)
		s.countPlays = 0
	}
}

type soundManager struct {
	sounds map[string]*sound
}

func newSoundManager() *soundManager {

	sm := soundManager{
		sounds: make(map[string]*sound)}
	return &sm
}

func (sm *soundManager) playSound(name string) {

	sound := sm.sounds[name]
	if sound == nil {
		fmt.Printf("Sound name %s has not found\n", name)
	} else {
		sound.play()
	}
}

func (sm *soundManager) loadSound(name, fileName string) {

	sound := newSound(fileName)
	sm.sounds[name] = sound
}
