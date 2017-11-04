package main

import (
	"fmt"
	"math"
	"math/rand"
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

func (p person) String() string {
	return fmt.Sprintf("%v is a professional %v and has been performing for %v years. Height: %vm. Weight: %vkg.", p.name, p.profession, p.timeOfPerformance, p.height, p.weight)
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
	person{"racedriver", "Fernando Alonso", 1.71, 68, 18},
	person{"racedriver", "Ayrton Senna", 1.75, 65, 10},
	person{"racedriver", "Jeff Gordon", 1.73, 68, 24},
	person{"racedriver", "Jimmie Johnson", 1.8, 75, 19},
	person{"racedriver", "Danica Patrick", 1.57, 45, 12},
	person{"racedriver", "Max Verstappen", 1.8, 67, 4},
	person{"racedriver", "Valtteri Bottas", 1.73, 70, 5},
	person{"racedriver", "Sergio Pérez", 1.73, 63, 6},
}

/*

var randomGuys = []person{
	person{"actor", "Daniel Radcliffe", 1.65, 74, 18},
	person{"singer", "Ricardo Arjona", 1.98, 92, 33},
	person{"singer", "Florence Welch", 1.75, 57, 10},
	person{"",,,},
	person{"",,,},
	person{"",,,},
	person{"",,,},
	person{"",,,},
	person{"",,,},
	person{"",,,},
}

*/

func shufflePersons(arr []person) []person {
	dest := make([]person, len(arr))
	perm := rand.Perm(len(arr))
	for i, v := range perm {
		dest[v] = arr[i]
	}
	return dest
}

func train(desiredType string, trainingData []person, weights []float64) []float64 {
	targetSuccesses := len(trainingData)
	iterations := 0
	learnRate := 0.5
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

		fmt.Printf("Finished run %v, Hit rate: %v/%v (%v%%)\n", iterations, successes, targetSuccesses, 100*successes/targetSuccesses)

		if successes == targetSuccesses || iterations > 10000 {
			return weights
		}
	}
}

func main() {
	rand.Seed(int64(time.Now().Nanosecond()))

	trainingData := shufflePersons(append(basketballers, racedrivers...))

	weights := make([]float64, 0)

	for i := 0; i < 4; i++ {
		weights = append(weights, 2*rand.Float64()-1)
	}

	fmt.Printf("Starting training with weights: %v\n\n", weights)

	newWeights := train("basketballer", trainingData, weights)

	fmt.Printf("\n\nFinished! The new weights are: %v", newWeights)
}
