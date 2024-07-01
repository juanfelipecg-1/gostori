package flags

import (
	"flag"
	"log"
)

var (
	defaultPath  = "./csv-example/txns.csv"
	defaultEmail = "default@gmail.com"
)

var (
	pathFlagMessage  = "Path to the transaction file"
	emailFlagMessage = "Email address to send the summary to"
)

func ParseFlags() (string, string) {
	// Define command-line flags
	filePath := flag.String("file", defaultPath, pathFlagMessage)
	email := flag.String("email", defaultEmail, emailFlagMessage)

	flag.Parse()

	// Validate command-line flags
	if *filePath == "" {
		log.Fatal("file path is required")
	}
	if *email == "" {
		log.Fatal("email is required")
	}

	return *filePath, *email
}
