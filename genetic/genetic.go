package main

import (
	"clase-sistemas-inteligentes/utilities"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func main() {
	rand.Seed(int64(time.Now().Nanosecond()))

	function := func(genotipe string) float64 {
		value, _ := strconv.ParseUint(genotipe, 2, 0)

		return float64(value)
	}

	population := generateInitialPopulation(50, 8)

	aptitudes := calculateAptitude(population, function)

	printAverageAptitude(aptitudes)

	for i := 0; i < 5000; i++ {

		aptitudes := calculateAptitude(population, function)

		population = generateOffspring(population, aptitudes)

	}

	aptitudes = calculateAptitude(population, function)

	printAverageAptitude(aptitudes)

	utilities.PrettyPrint(population)

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

func calculateAptitude(population []string, function func(individual string) float64) []float64 {
	aptitude := make([]float64, len(population))

	for i := 0; i < len(aptitude); i++ {
		aptitude[i] = function(population[i])
	}

	return aptitude
}

func generateOffspring(population []string, aptitudes []float64) []string {
	aptitudeSum := float64(0)
	parents := make([]string, len(population))
	couples := make([]string, len(population))
	children := make([]string, len(population))
	permutation := rand.Perm(len(population))

	for _, aptitude := range aptitudes {
		aptitudeSum += aptitude
	}

	// Choose parents
	for i := 0; i < len(population); i++ {
		rouletteNumber := rand.Float64() * aptitudeSum
		runningSum := float64(0)

		for j := 0; j < len(aptitudes); j++ {
			runningSum += aptitudes[j]
			if rouletteNumber < runningSum {
				parents[i] = population[j]
				break
			}
		}
	}

	// Make couples
	for i, v := range permutation {
		couples[v] = parents[i]
	}

	// Make offstring
	for i := 0; i < len(couples); i += 2 {
		crossPoint := rand.Intn(len(couples[i]) - 1)
		coupleChildren := [2]string{"", ""}
		for j := 0; j < len(couples[i]); j++ {
			if j > crossPoint {
				coupleChildren[0] += string(couples[i][j])
				coupleChildren[1] += string(couples[i+1][j])
			} else {
				coupleChildren[0] += string(couples[i+1][j])
				coupleChildren[1] += string(couples[i][j])
			}
		}
		children[i] = coupleChildren[0]
		children[i+1] = coupleChildren[1]
	}

	// Mutate

	for i := 0; i < len(children); i++ {
		if rand.Float32() < 0.05 {
			mutatedGene := rand.Intn(len(children[i]))
			if children[i][mutatedGene] == '1' {
				children[i] = replaceAtIndex(children[i], '0', mutatedGene)
			} else {
				children[i] = replaceAtIndex(children[i], '1', mutatedGene)
			}
		}
	}

	return children
}

func replaceAtIndex(str string, replacement rune, index int) string {
	return str[:index] + string(replacement) + str[index+1:]
}

func printAverageAptitude(aptitudes []float64) {
	aptitudeSum := float64(0)

	for _, aptitude := range aptitudes {
		aptitudeSum += aptitude
	}

	fmt.Printf("Average aptitude: %f\n", (aptitudeSum / float64(len(aptitudes))))
}
