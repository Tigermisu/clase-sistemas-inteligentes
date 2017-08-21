package main

import (
	"clase-sistemas-inteligentes/utilities"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// A GuessTree is a binary tree that also contains an attribute to identify it as a question or not
type GuessTree struct {
	Left       *GuessTree
	Value      string
	Right      *GuessTree
	IsQuestion bool
}

var guessTree GuessTree
var usingLocal = false

const jsonURL = "https://api.myjson.com/bins/a9101"

func main() {
	startGuessing()
}

func loadGuessTree() {
	treeData, err := utilities.GetData(jsonURL)
	if err != nil {
		// We may not have internet, fall back to local copy
		treeData, err = ioutil.ReadFile("tree.json")
		usingLocal = true // Keep track of local copy usage, just so we don't overwrite global state
		if err != nil {
			// No local copy, create one
			guessTree = createDefaultGuessTree()
			saveGuessTree()
		}
	} else {
		usingLocal = false
	}

	err = json.Unmarshal(treeData, &guessTree)
	if err != nil {
		// Malformed JSON, fall back to local copy
		fmt.Println("Oops, I think my memory is gone.")
		guessTree = createDefaultGuessTree()
		saveGuessTree()
	}

}

func saveGuessTree() {
	if !usingLocal {
		utilities.PutJSONData(jsonURL, guessTree)
	}

	jsonData, err := json.Marshal(guessTree)
	utilities.AssertError(err)
	err = ioutil.WriteFile("tree.json", []byte(jsonData), 0644)
	utilities.AssertError(err)
}

func createDefaultGuessTree() GuessTree {
	var newTree GuessTree
	newTree.IsQuestion = true
	newTree.Value = "eats meat?"
	newTree.Left = &GuessTree{nil, "Cow", nil, false}
	newTree.Right = &GuessTree{nil, "T-Rex", nil, false}

	return newTree
}

func askForYesOrNo(longPrompt bool) bool {
	for {
		if longPrompt {
			fmt.Print("Please enter 'y'es or 'n'o (y/n): ")
		} else {
			fmt.Print("(y/n): ")
		}
		answer := utilities.GetConsoleInput()
		if answer != "y" && answer != "n" {
			fmt.Print("Sorry, invalid response. ")
		} else {
			return answer == "y"
		}
	}
}

func guess() bool {
	currentNode := &guessTree
	for {
		if currentNode.IsQuestion {
			fmt.Printf("Your animal... %s? ", currentNode.Value)
			if askForYesOrNo(false) {
				currentNode = currentNode.Right
			} else {
				currentNode = currentNode.Left
			}
		} else {
			fmt.Printf("Aha! Your animal is: %s.\n", currentNode.Value)
			fmt.Print("Did I win? ")

			if askForYesOrNo(false) {
				fmt.Printf("Excellent! One more chocolate for me!")
				return true
			}

			fmt.Print("Aww, bummer. Can you please tell me the name of your animal?\nEnter name: ")
			name := utilities.GetConsoleInput()
			fmt.Printf("Okay, can you please now tell me one distinctive characteristic of this animal? (ie. 'has blue feathers', 'runs very fast', 'flies for long periods of time'...)\nEnter characteristic: ")
			characteristic := utilities.GetConsoleInput()
			addNewGuessOption(currentNode, name, characteristic)
			fmt.Printf("Alright, I'll remember that one for the next time.")
			return false
		}
	}
}

func addNewGuessOption(lastNode *GuessTree, newName string, newCharacteristic string) {
	oldName := lastNode.Value
	lastNode.IsQuestion = true
	lastNode.Value = newCharacteristic
	lastNode.Left = &GuessTree{nil, oldName, nil, false}
	lastNode.Right = &GuessTree{nil, newName, nil, false}
	saveGuessTree()
}

func startGuessing() {
	fmt.Println("Hello! My name is Chris Intelligent Guesser 3000 version Exodia Prime RX.")
	fmt.Println("The game is simple. You think of an animal. I ask a series of yes/no questions to try and guess it.")
	fmt.Println("If I fail to guess it, you win! But if I guess it, you owe me a chocolate.")
	fmt.Println("Are you ready?")
	for !askForYesOrNo(true) {
		fmt.Println("Oh, come on! I'll ask again.")
	}
	fmt.Println("\nGreat! Let's begin!")
	playAgain := true
	correctGuesses := 0
	for playAgain {
		loadGuessTree()
		won := guess()
		if won {
			correctGuesses++
		}
		fmt.Print("\n\nDo you want to play another round? ")
		playAgain = askForYesOrNo(false)
	}
	if correctGuesses > 0 {
		if correctGuesses == 1 {
			fmt.Printf("Okay, just remember, you owe me %d chocolate ;)", correctGuesses)
		} else {
			fmt.Printf("Okay, just remember, you owe me %d chocolates ;)", correctGuesses)
		}
	} else {
		fmt.Println("Fine, I'll get my own chocolates.")
	}
}
