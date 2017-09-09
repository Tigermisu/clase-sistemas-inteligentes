package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"
)

var distances = [7][7]float64{
	{0, 1802, 1180, 880, 1460, 1394, 809},   // Chihuahua
	{1792, 0, 1829, 1339, 418, 2897, 999},   // Veracruz
	{1244, 1835, 0, 1745, 1820, 1274, 918},  // Dallas
	{879, 1392, 1845, 0, 1019, 2462, 842},   // Mazatlan
	{1424, 395, 1809, 1022, 0, 27476, 910},  // CDMX
	{1407, 2894, 1285, 2350, 2746, 0, 1913}, // Denver
	{804, 999, 915, 839, 913, 1913, 0},      // Monterrey
}

func main() {
	rand.Seed(int64(time.Now().Nanosecond()))

	iterations := 0

	cities := generateInitialCities(200)

	distanceFunction := func(order [6]int) float64 {
		distance := distances[0][order[0]]
		for i := 0; i < len(order)-1; i++ {
			distance += distances[order[i]][order[i+1]]
		}
		return 1000000 / (distance + distances[order[len(order)-1]][0])
	}

	aptitudes := calculatePathAptitude(cities, distanceFunction)

	oldAvg := math.Inf(1)

	hits := 0

	for {
		iterations++

		aptitudes := calculatePathAptitude(cities, distanceFunction)

		cities = generateNewPath(cities, aptitudes)

		aptitudes = calculatePathAptitude(cities, distanceFunction)

		newAvg := calculateAverageAptitude(aptitudes)

		if math.Abs(newAvg-oldAvg) < 0.5 {
			hits++
			if hits > 9 {
				break
			}
		} else {
			hits = 0
		}
		oldAvg = newAvg

	}

	aptitudes = calculatePathAptitude(cities, distanceFunction)

	average := calculateAverageAptitude(aptitudes)

	//fmt.Println(cities)

	fmt.Printf("Average aptitude: %f, Iterations: %d\n", average, iterations)

	printBestCity(cities, distanceFunction)

}

func printBestCity(cities [][6]int, function func([6]int) float64) {
	var bestCity = [6]int{}
	maxAptitude := 0.0

	for i := 0; i < len(cities); i++ {
		newAptitude := function(cities[i])
		if newAptitude > maxAptitude {
			bestCity = cities[i]
			maxAptitude = newAptitude
		}
	}

	fmt.Println("Best City:")
	fmt.Print(bestCity)
	fmt.Printf(" %f", maxAptitude)
}
func generateInitialCities(popSize int) [][6]int {
	population := make([][6]int, popSize)

	for i := 0; i < popSize; i++ {
		cities := [6]int{1, 2, 3, 4, 5, 6}
		randomCities := [6]int{}
		permutation := rand.Perm(6)

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

func calculatePathAptitude(cities [][6]int, function func([6]int) float64) []float64 {
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

func generateNewPath(population [][6]int, aptitudes []float64) [][6]int {
	aptitudeSum := 0.0
	parents := make([][6]int, len(population))
	couples := make([][6]int, len(population))
	children := make([][6]int, len(population))
	permutation := rand.Perm(len(population))

	for _, aptitude := range aptitudes {
		aptitudeSum += aptitude
	}

	// Choose parents
	for i := 0; i < len(population); i++ {
		rouletteNumber := rand.Float64() * aptitudeSum
		runningSum := 0.0

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
		coupleChildren := [2][6]int{}

		for j := 0; j < len(couples[i]); j++ {
			if j <= crossPoint {
				coupleChildren[0][j] = couples[i][j]
				coupleChildren[1][j] = couples[i+1][j]
			} else {
				for k := 0; k < len(couples[i]); k++ {
					if !contains(coupleChildren[0], couples[i+1][k]) {
						coupleChildren[0][j] = couples[i+1][k]
						break
					}
				}
				for k := 0; k < len(couples[i]); k++ {
					if !contains(coupleChildren[1], couples[i][k]) {
						coupleChildren[1][j] = couples[i][k]
						break
					}
				}
			}
		}

		children[i] = coupleChildren[0]
		children[i+1] = coupleChildren[1]
	}

	// Mutate
	for i := 0; i < len(children); i++ {
		if rand.Float32() < 0.05 {
			children[i] = swapRandom(children[i])
		}
	}

	return children
}

func contains(array [6]int, value int) bool {
	for i := 0; i < len(array); i++ {
		if array[i] == value {
			return true
		}
	}
	return false
}

func swapRandom(arr [6]int) [6]int {
	permutation := rand.Perm(4)

	a := arr[permutation[0]]
	arr[permutation[0]] = arr[permutation[1]]
	arr[permutation[1]] = a

	return arr
}

func generateOffspring(population []string, aptitudes []float64) []string {
	aptitudeSum := 0.0
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
		runningSum := 0.0

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
	aptitudeSum := 0.0

	for _, aptitude := range aptitudes {
		aptitudeSum += aptitude
	}

	return aptitudeSum / float64(len(aptitudes))
}
