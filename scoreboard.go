package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	border  = 4
	shieldW = 100
	shieldH = 12
)

type scoreboard struct {
	texts       map[string]*text
	p1Points    int32
	p1Shield    int8
	hiScore     int32
	secLinePosY float64
}

func newScoreboard(renderer *sdl.Renderer) (s scoreboard) {

	s.p1Points = 0

	s.texts = make(map[string]*text)
	textP1 := newText(renderer, border, border, "Player1 0")
	s.texts["pointsP1"] = &textP1
	textHSLable := newText(renderer, screenWidth/2, border, "Hi Score")
	textHSLable.x = screenWidth/2 - float64(textHSLable.texW*scale/2)
	s.texts["hiScoreLable"] = &textHSLable

	s.secLinePosY = float64(textHSLable.texH*scale + border*2)
	textHS := newText(renderer, screenWidth/2, s.secLinePosY, "0")
	textHS.x = screenWidth/2 - float64(textHS.texW*scale/2)
	s.texts["hiScore"] = &textHS
	textP1Lives := newText(renderer, shieldW+border, s.secLinePosY, "0")
	s.texts["livesP1"] = &textP1Lives

	return s
}

func (s *scoreboard) draw(renderer *sdl.Renderer) {

	s.texts["pointsP1"].draw(renderer)
	s.texts["livesP1"].draw(renderer)
	s.texts["hiScoreLable"].draw(renderer)
	s.texts["hiScore"].draw(renderer)

	w := shieldW / playerMaxShield * s.p1Shield
	fillRect := sdl.Rect{X: border, Y: int32(s.secLinePosY), W: int32(w - border), H: shieldH * scale}
	renderer.SetDrawColor(92, 70, 140, 255)
	renderer.FillRect(&fillRect)
	renderer.SetDrawColor(200, 200, 0, 255)
	drawRect := sdl.Rect{X: border, Y: int32(s.secLinePosY), W: shieldW - border, H: shieldH * scale}
	renderer.DrawRect(&drawRect)
	drawRectIn := sdl.Rect{X: border + 1, Y: int32(s.secLinePosY) + 1, W: shieldW - border - 2, H: shieldH*scale - 2}
	renderer.DrawRect(&drawRectIn)
}

func (s *scoreboard) resetPointsP1() {

	s.p1Points = 0
	s.texts["pointsP1"].setValue(fmt.Sprintf("Player1 %d", s.p1Points))
}

func (s *scoreboard) incrementPointsP1(qtd int32) {

	s.p1Points += qtd
	s.texts["pointsP1"].setValue(fmt.Sprintf("Player1 %d", s.p1Points))
	s.updateHighScore()
}

func (s *scoreboard) setLivesP1(qtd int8) {

	s.texts["livesP1"].setValue(fmt.Sprintf("%d", qtd))
}

func (s *scoreboard) setShieldP1(qtd int8) {

	s.p1Shield = qtd
}

func (s *scoreboard) updateHighScore() {

	if s.p1Points > s.hiScore {
		s.hiScore = s.p1Points
		hs := s.texts["hiScore"]
		hs.setValue(fmt.Sprintf("%d", s.hiScore))
		hs.x = screenWidth/2 - float64(hs.texW*scale/2)
	}
}

func (s *scoreboard) update() {

	// do nothing
}
