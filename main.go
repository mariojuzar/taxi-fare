package main

import (
	"bufio"
	"github.com/mariojuzar/taxi-fare/fare"
	"log"
	"os"
	"strings"
)

func main() {
	// initialize taxi fare
	taxiFare := fare.NewTaxiFare()

	var (
		err         error
		input       string
		currentFare int64
	)

	log.Println("Please input distance meter records:")
	reader := bufio.NewReader(os.Stdin)

	for err == nil && input != "finished" {
		input, err = reader.ReadString('\n')
		input = strings.TrimSuffix(input, "\n")
		currentFare, err = taxiFare.CalculateFareMeter(input)
		if err != nil {
			log.Printf("Error :%s", err.Error())
			break
		}
		log.Printf("Current fare: %d\n", currentFare)
	}

	if err != nil {
		os.Exit(1)
	}
}
