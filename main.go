package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type GameState struct {
	words       []string
	choosenWord string
	showWord    string
}

type Request struct {
	Letter string `json:"letter"`
}

var gameState GameState

func getWord(c echo.Context) error {
	wrd := fmt.Sprintf("Word: %s", gameState.showWord)
	return c.String(http.StatusOK, wrd)
}

func insertWord(c echo.Context) error {
	var req Request
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "Bad request")
	}

	won, message := processLetter(req.Letter)
	if won {
		return c.String(http.StatusOK, message) // Retorna mensagem de vitória
	}

	wrd := fmt.Sprintf("Word: %s", gameState.showWord)
	return c.String(http.StatusOK, wrd)
}

func newGame(c echo.Context) error {
	gameState.showWord = ""

	initGame()

	wrd := fmt.Sprintf("Word: %s", gameState.showWord)
	return c.String(http.StatusOK, wrd)
}

func main() {
	initGame()

	e := echo.New()
	e.GET("/word", getWord)
	e.GET("/word/new", newGame)
	e.POST("/word", insertWord)

	e.Logger.Fatal(e.Start(":1323"))
}

func initGame() {
	gameState.words = []string{"peteca", "luz", "janela", "abacaxi", "suquinho", "abacate", "bola", "desenho", "paralelepipedo"}

	// Gera indice aleatório
	reader := rand.New(rand.NewSource(time.Now().UnixNano()))
	index := reader.Intn(len(gameState.words))

	gameState.choosenWord = gameState.words[index]

	for range gameState.choosenWord {
		gameState.showWord += "_"
	}
}

func processLetter(letter string) (bool, string) {
	letter = strings.ToLower(letter)
	if strings.Contains(gameState.choosenWord, letter) {
		for i, c := range gameState.choosenWord {
			if string(c) == letter {
				gameState.showWord = gameState.showWord[:i] + letter + gameState.showWord[i+1:]
			}
		}
	}

	if !strings.Contains(gameState.showWord, "_") {
		txt := "Parabéns! Você ganhou!\n"
		txt += "A palavra era " + gameState.choosenWord
		return true, txt
	}

	return false, ""
}
