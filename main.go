package main

import (
	"log"
	"github.com/claudioontheweb/go-identicon/identicon"
)

func main() {

	title := "test"

	data := []byte(title)
	identiconElem := identicon.HashInput(data)

	identiconElem = identicon.Pipe(identiconElem, identicon.PickColor, identicon.BuildGrid, identicon.FilterOddSquares, identicon.BuildPixelMap)

	if err := identicon.DrawRectangle(identiconElem); err != nil {
		log.Fatalln(err)
	}

}