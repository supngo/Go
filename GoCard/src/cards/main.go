package main

import "fmt"

// var card string

func main() {
	// var card string = "Ace of Spades"
	// card := "Ace of Spades"
	// card = "Five of Diamonds"

	// card := newCard()
	// fmt.Println(card)

	// cards := deck{newCard(), "Ace of Diamonds"}
	// cards = append(cards, "Six of Hearts")

	// for i, iter := range cards {
	// 	fmt.Println(i, iter)
	// }
	// fmt.Println(cards)

	cards := newDeck()
	fmt.Println(cards.toString())
	cards.saveToFile("thom_deck")

	newDeck := newDeckFromFile("thom_deck")
	newDeck.shuffle()
	newDeck.print()

	// hand, remainingDeck := deal(cards, 5)
	// hand.print()
	// remainingDeck.print()
}
