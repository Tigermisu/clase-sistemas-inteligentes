package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

/*
Stats:

Height
Weight
Fitness

*/

type person struct {
	profession        string
	name              string
	height            float64
	weight            float64
	timeOfPerformance float64
}

type statistics struct {
	key      float64
	mean     float64
	variance float64
	stddev   float64
}

func (p person) String() string {
	return fmt.Sprintf("%v is a professional %v and has been performing for %v years. Height: %vm. Weight: %vkg.", p.name, p.profession, p.timeOfPerformance, p.height, p.weight)
}

func (s statistics) String() string {
	return fmt.Sprintf("%v:\n\tMean: %v.\n\tVariance: %v\n\tStd. Dev.: %v\n\n", s.key, s.mean, s.variance, s.stddev)
}

var basketballers = []person{
	person{"basketballer", "Shaquille O'Neal", 2.16, 147, 19},
	person{"basketballer", "LeBron James", 2.03, 113, 14},
	person{"basketballer", "Stephen Curry", 1.91, 86, 12},
	person{"basketballer", "Kobe Bryant", 1.98, 96, 20},
	person{"basketballer", "Dirk Nowitzki", 2.13, 111, 24},
	person{"basketballer", "Yao Ming", 2.29, 141, 14},
	person{"basketballer", "Matthew Dellavedova", 1.93, 90, 11},
	person{"basketballer", "Tim Duncan", 2.11, 113, 19},
	person{"basketballer", "Kawhi Leonard", 2.01, 104, 10},
	person{"basketballer", "Michael Jordan", 1.98, 98, 19},
}

var racedrivers = []person{
	person{"racedriver", "Sebastian Vettel", 1.75, 62, 12},
	person{"racedriver", "Lewis Hamilton", 1.74, 68, 10},
	person{"racedriver", "Fernando Alonso", 1.71, 68, 14},
	person{"racedriver", "Ayrton Senna", 1.75, 65, 10},
	person{"racedriver", "Jeff Gordon", 1.73, 68, 14},
	person{"racedriver", "Jimmie Johnson", 1.8, 75, 9},
	person{"racedriver", "Danica Patrick", 1.57, 45, 12},
	person{"racedriver", "Max Verstappen", 1.8, 67, 4},
	person{"racedriver", "Valtteri Bottas", 1.73, 70, 5},
	person{"racedriver", "Sergio Pérez", 1.73, 63, 6},
}

/*

var oldpeople = []person{
	person{"oldperson", "Mack Oldperson", 1.82, 80, 52},
	person{"oldperson", "Gramps McGramps", 1.75, 64, 34},
	person{"oldperson", "George Dinosaur", 1.98, 120, 42},
	person{"oldperson", "Granny Gran", 1.53, 45, 45},
	person{"oldperson", "Miss Matusalen", 1.45, 42, 70},
	person{"oldperson", "Speedy Gonzalez", 1.46, 56, 35},
}

var randomGuys = []person{
	person{"actor", "Daniel Radcliffe", 1.65, 74, 18},
	person{"singer", "Ricardo Arjona", 1.98, 92, 33},
	person{"singer", "Florence Welch", 1.75, 57, 10},
	person{"kpop artist", "Solar", 1.61, 43, 4},
	person{"kpop artist", "Moon Byul", 1.63, 45, 4},
	person{"kpop artist", "Whee In", 1.59, 43, 4},
	person{"kpop artist", "Hwa Sa", 1.6, 44, 4},
}

*/

func main() {
	args := os.Args
	fmt.Println(`Welcome to Chris' simple classifier.
	Run without arguments to train the perceptron and generate new weights.
	Run with one argument (sample size) to optimize learn rate and generate statistics.
	Run with 3 arguments (Height in meters, weight in kilos and expertise in years) to attempt a classification.`)

	if len(args) == 4 {
		height, _ := strconv.ParseFloat(args[1], 64)
		weight, _ := strconv.ParseFloat(args[2], 64)
		expertise, _ := strconv.ParseFloat(args[3], 64)
		classify(person{"Unknown", "Your person", height, weight, expertise}, false)
	} else {
		rand.Seed(int64(time.Now().Nanosecond()))

		trainingData := shufflePersons(append(basketballers, racedrivers...))

		basketWeights := make([]float64, 0)
		racedriverWeights := make([]float64, 0)

		for i := 0; i < 4; i++ {
			basketWeights = append(basketWeights, 2*rand.Float64()-1)
		}

		if len(args) == 2 {
			sampleSize, _ := strconv.Atoi(args[1])
			optimizeLearnRate("basketballer", trainingData, basketWeights, sampleSize)
		} else {
			fmt.Println("\n\nReady to train! The training data consists of:")

			for _, p := range trainingData {
				fmt.Println(p)
			}

			fmt.Printf("\n\nStarting basketball training with weights:\n%v\n\n", basketWeights)

			train("basketballer", trainingData, basketWeights, 0.5)

			fmt.Printf("Finished! The new basketball weights are:\n%v\n", basketWeights)

			for i := 0; i < 4; i++ {
				racedriverWeights = append(racedriverWeights, 2*rand.Float64()-1)
			}

			fmt.Printf("\n\nStarting racedriver training with weights:\n%v\n\n", racedriverWeights)

			train("racedriver", trainingData, racedriverWeights, 0.7)

			fmt.Printf("Finished! The new racedriver weights are:\n%v\n\n", racedriverWeights)

			fmt.Println("Writing weights to trained-weights.txt")

			f, err := os.Create("trained-weights.txt")
			if err != nil {
				panic(err)
			}

			defer f.Close()

			f.WriteString("[ Basketball ]\n")

			for i := 0; i < 4; i++ {
				f.WriteString(fmt.Sprintf("%v\n", basketWeights[i]))
			}

			f.WriteString("[ Racecars ]\n")

			for i := 0; i < 4; i++ {
				f.WriteString(fmt.Sprintf("%v\n", racedriverWeights[i]))
			}

			f.Sync()

			fmt.Println("\nDone!")
		}
	}

}

func classify(p person, shouldPanic bool) {
	var fileLine string

	basketWeights := make([]float64, 4)
	racedriverWeights := make([]float64, 4)

	file, err := os.Open("trained-weights.txt")
	defer file.Close()

	if err != nil {
		if !shouldPanic {
			os.Args = make([]string, 0)
			main()
			classify(p, true)
		} else {
			panic(err)
		}
	}

	reader := bufio.NewReader(file)
	fileLine, err = reader.ReadString('\n')
	fmt.Println("\n\nRetrieving stored weights.\n\nReading " + fileLine)

	for i := 0; i < 4; i++ {
		fileLine, err = reader.ReadString('\n')
		cleanLine := strings.TrimSuffix(fileLine, "\n")
		basketWeights[i], _ = strconv.ParseFloat(cleanLine, 64)
	}

	fileLine, err = reader.ReadString('\n')
	fmt.Println("Reading " + fileLine)

	for i := 0; i < 4; i++ {
		fileLine, err = reader.ReadString('\n')
		cleanLine := strings.TrimSuffix(fileLine, "\n")
		racedriverWeights[i], _ = strconv.ParseFloat(cleanLine, 64)
	}

	fmt.Printf("Retrieved basketball weights:\n%v\n\n", basketWeights)
	fmt.Printf("Retrieved racedriver weights:\n%v\n\n", racedriverWeights)

	basketSum := basketWeights[0] + basketWeights[1]*p.height + basketWeights[2]*p.weight + basketWeights[3]*p.timeOfPerformance
	raceSum := racedriverWeights[0] + racedriverWeights[1]*p.height + racedriverWeights[2]*p.weight + racedriverWeights[3]*p.timeOfPerformance

	if basketSum > 0 {
		fmt.Printf("%v could be a basket baller!\n", p.name)
	} else {
		fmt.Printf("%v could not be a basket baller :(\n", p.name)
	}

	if raceSum > 0 {
		fmt.Printf("%v could be a race driver!\n", p.name)
	} else {
		fmt.Printf("%v could not be a race driver :(\n", p.name)
	}
}

func train(desiredType string, trainingData []person, weights []float64, learnRate float64) ([]float64, int) {
	targetSuccesses := len(trainingData)
	iterations := 0
	for {
		iterations++
		successes := 0

		for _, p := range trainingData {
			sum := weights[0] + weights[1]*p.height + weights[2]*p.weight + weights[3]*p.timeOfPerformance

			if (sum > 0 && p.profession != desiredType) || (sum <= 0 && p.profession == desiredType) {
				// We got it wrong :c
				var expectedY float64

				if p.profession == desiredType {
					expectedY = 1.0
				} else {
					expectedY = -1.0
				}

				weights[0] += learnRate * expectedY
				weights[1] += learnRate * expectedY * weights[1]
				weights[2] += learnRate * expectedY * weights[2]
				weights[3] += learnRate * expectedY * weights[3]

				// If a weight is absolutely more than 10 or devilishly small, we reset it to a new random. This prevents diverging
				for i := 0; i < 4; i++ {
					if math.Abs(weights[i]) >= 10 || math.Abs(weights[i]) <= 1.0e-10 {
						weights[i] = 2*rand.Float64() - 1
					}
				}

			} else {
				// All good :D
				successes++
			}

		}

		if successes == targetSuccesses || iterations > 10000 {
			//fmt.Printf("Finished run for %v. Iterations: %v, Hit rate: %v/%v (%v%%)\n", desiredType, iterations, successes, targetSuccesses, 100*successes/targetSuccesses)
			return weights, iterations
		}
	}
}

func optimizeLearnRate(desiredType string, trainingData []person, weights []float64, sampleSize int) {
	fmt.Printf("\nOptimizing with sample size of %v...\n\n", sampleSize)
	stats := make([]statistics, 0)
	loopCount := 0
	start := 0.0005
	limit := 3.5
	total := limit / start
	for learnRate := start; learnRate <= limit; learnRate += start {
		iterationSamples := make([]int, 0)
		for run := 0; run < sampleSize; run++ {

			_, iterations := train(desiredType, trainingData, weights, learnRate)

			iterationSamples = append(iterationSamples, iterations)

			for i := 0; i < 4; i++ {
				weights[i] = 2*rand.Float64() - 1
			}
		}

		mean := getMean(iterationSamples)
		variance := getVariance(iterationSamples, mean)

		stats = append(stats, statistics{learnRate, mean, variance, math.Sqrt(variance)})
		loopCount++

		if loopCount%100 == 0 {
			fmt.Printf("Progress: %v%%\r", 100*loopCount/int(total))
		}
	}

	fmt.Println("Sorting...")

	sort.Slice(stats, func(i, j int) bool { return stats[i].mean < stats[j].mean })

	fmt.Println("Writing to stats.txt")

	f, err := os.Create("stats.txt")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	for _, s := range stats {
		f.WriteString(s.String())
		f.Sync()
	}

	fmt.Println("Writing to chartable.csv")

	f, err = os.Create("chartable.csv")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	for _, s := range stats {
		f.WriteString(fmt.Sprintf("%v,%v,%v,%v\n", s.key, s.mean, s.variance, s.stddev))
		f.Sync()
	}

	fmt.Println("Done!")
}

func shufflePersons(arr []person) []person {
	dest := make([]person, len(arr))
	perm := rand.Perm(len(arr))
	for i, v := range perm {
		dest[v] = arr[i]
	}
	return dest
}

func getMean(data []int) float64 {
	sum := 0

	for _, v := range data {
		sum += v
	}

	return float64(sum) / float64(len(data))
}

func getVariance(data []int, mean float64) float64 {
	sum := 0.0

	for _, v := range data {
		sum += math.Pow(float64(v)-mean, 2)
	}

	return sum / float64(len(data)-1)
}
