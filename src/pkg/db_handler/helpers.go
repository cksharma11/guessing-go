package dbhandler

import (
	"math/rand"
	"strconv"
)

func getWinner(guesses map[string]string) string {
	if len(guesses) == 0 {
		return "Anonymous"
	}

	guessMap := getGuessMap(guesses)
	maxGuess := 0
	var maxGuessBy []string
	for guess, guessedBy := range guessMap {
		if maxGuess < guess {
			maxGuessBy = guessedBy
		}
	}

	winnerIndex := rand.Intn(len(maxGuessBy))
	return maxGuessBy[winnerIndex]
}

func getGuessMap(guesses map[string]string) map[int][]string {
	guessMap := make(map[int][]string)
	for _, key := range guesses {
		guess, _ := strconv.Atoi(key)
		guessMap[guess] = []string{}
	}
	for value, key := range guesses {
		guess, _ := strconv.Atoi(key)
		guessMap[guess] = append(guessMap[guess], value)
	}
	return guessMap
}
