package main

import (
	"clase-sistemas-inteligentes/utilities"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"
)

var distances = [5][5]float64{
	{0, 1802, 1180, 880, 1460},  // Chihuahua
	{1792, 0, 1829, 1339, 418},  // Veracruz
	{1244, 1835, 0, 1745, 1820}, // Dallas
	{879, 1392, 1845, 0, 1019},  // Mazatlan
	{1424, 395, 1809, 1022, 0},  // CDMX
}

func main() {
	rand.Seed(int64(time.Now().Nanosecond()))

	iterations := 0

	cities := generateInitialCities(6)

	distanceFunction := func(order [4]int) float64 {
		distance := distances[0][order[0]]
		for i := 0; i < 3; i++ {
			distance += distances[order[i]][order[i+1]]
		}
		return distance + distances[order[3]][0]
	}

	aptitudes := calculatePathAptitude(cities, distanceFunction)

	oldAvg := math.Inf(1)

	hits := 0

	for iterations < 500000 {
		iterations++

		aptitudes := calculatePathAptitude(cities, distanceFunction)

		cities = generateNewPath(cities, aptitudes)

		aptitudes = calculatePathAptitude(cities, distanceFunction)

		newAvg := calculateAverageAptitude(aptitudes)

		if math.Abs(newAvg-oldAvg) < 10 {
			hits++
			if hits > 9 {
				break
			}
		} else {
			hits = 0
		}
		oldAvg = newAvg

	}

	/*

		function := func(genotipe string) float64 {
			value, _ := strconv.ParseUint(genotipe, 2, 0)
			if value >= 127 {
				value = 255 - value
			}
			return float64(value * 2)
		}

		population := generateInitialPopulation(20, 8)

		aptitudes := calculateAptitude(population, function)

		calculateAverageAptitude(aptitudes)


		for {

			aptitudes := calculateAptitude(population, function)

			if calculateAverageAptitude(aptitudes) >= 254 {
				break
			}

			population = generateOffspring(population, aptitudes)
			iterations++

		}

	*/

	aptitudes = calculatePathAptitude(cities, distanceFunction)

	average := calculateAverageAptitude(aptitudes)

	fmt.Printf("Average aptitude: %f\n, Iterations: %d\n\n", average, iterations)

	utilities.PrettyPrint(cities)

}

func generateInitialCities(popSize int) [][4]int {
	population := make([][4]int, popSize)

	for i := 0; i < popSize; i++ {
		cities := [4]int{1, 2, 3, 4}
		randomCities := [4]int{}
		permutation := rand.Perm(4)

		for i, v := range permutation {
			randomCities[v] = cities[i]
		}

		population[i] = randomCities
	}

	return population
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

func calculatePathAptitude(cities [][4]int, function func([4]int) float64) []float64 {
	aptitude := make([]float64, len(cities))

	for i := 0; i < len(aptitude); i++ {
		aptitude[i] = function(cities[i])
	}

	return aptitude
}

func calculateAptitude(population []string, function func(individual string) float64) []float64 {
	aptitude := make([]float64, len(population))

	for i := 0; i < len(aptitude); i++ {
		aptitude[i] = function(population[i])
	}

	return aptitude
}

func generateNewPath(population [][4]int, aptitudes []float64) [][4]int {
	aptitudeSum := float64(0)
	parents := make([][4]int, len(population))
	couples := make([][4]int, len(population))
	children := make([][4]int, len(population))
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

	// Make offspring
	for i := 0; i < len(couples); i += 2 {
		crossPoint := rand.Intn(len(couples[i]) - 1)
		coupleChildren := [2][4]int{}

		children[i] = coupleChildren[0]
		children[i+1] = coupleChildren[1]
	}

	fmt.Println(children)

	// Mutate
	for i := 0; i < len(children); i++ {
		if rand.Float32() < 0.05 {
			children[i] = swapRandom(children[i])
		}
	}

	return children
}

func swapRandom(arr [4]int) [4]int {
	permutation := rand.Perm(4)

	a := arr[permutation[0]]
	arr[permutation[0]] = arr[permutation[1]]
	arr[permutation[1]] = a

	return arr
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

func calculateAverageAptitude(aptitudes []float64) float64 {
	aptitudeSum := float64(0)

	for _, aptitude := range aptitudes {
		aptitudeSum += aptitude
	}

	return aptitudeSum / float64(len(aptitudes))
}
