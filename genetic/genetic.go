package main

import (
	"clase-sistemas-inteligentes/utilities"
	"math/rand"
	"strconv"
	"time"
)

// A GuessTree is a binary tree that also contains an attribute to identify it as a question or not
type GuessTree struct {
	Left       *GuessTree
	Value      string
	Right      *GuessTree
	IsQuestion bool
}

const jsonURL = "https://api.myjson.com/bins/a9101"

func main() {
	rand.Seed(int64(time.Now().Nanosecond()))

	function := func(genotipe string) float32 {
		value, _ := strconv.ParseUint(genotipe, 2, 0)

		return float32(value)
	}

	population := generateInitialPopulation(12, 5)
	utilities.PrettyPrint(population)

	aptitudes := calculateAptitude(population, function)

	utilities.PrettyPrint(aptitudes)

}

func generateInitialPopulation(populationSize, dnaLength int) []string {
	population := make([]string, populationSize)

	for i := 0; i < populationSize; i++ {
		population[i] = generateRandomString(dnaLength)
	}

	return population
}

func generateRandomString(length int) string {
	str := ""

	for i := 0; i < length; i++ {
		str += strconv.Itoa(rand.Intn(2))
	}

	return str
}

func calculateAptitude(population []string, function func(individual string) float32) []float32 {
	aptitude := make([]float32, len(population))

	for i := 0; i < len(aptitude); i++ {
		aptitude[i] = function(population[i])
	}

	return aptitude
}
