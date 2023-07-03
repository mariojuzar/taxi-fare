package main

import (
	"bufio"
	"github.com/mariojuzar/taxi-fare/fare"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"strings"
)

func main() {
	// set log format to JSON
	logrus.SetFormatter(&logrus.JSONFormatter{})

	taxiFare := fare.NewTaxiFare()

	log.Println("Please input distance meter records:")
	reader := bufio.NewReader(os.Stdin)

	var finalFare int64

	for {
		input, err := reader.ReadString('\n')
		input = strings.TrimSuffix(input, "\n")

		if input == "end" {
			break
		}

		currentFare, err := taxiFare.Calculate(strings.TrimSpace(input))
		if err != nil {
			logrus.Fatalf("Error : %s", err.Error())
		}
		log.Printf("Current fare: %d\n", currentFare)
		finalFare = currentFare
	}

	log.Printf("Final fare: %d\n", finalFare)
	taxiFare.ShowRecords()
}
