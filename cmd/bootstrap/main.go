package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	_, debug := os.LookupEnv("LIFEOS_DEBUG")
	if debug {
		log.Println("Debugging requested")
		for _, pair := range os.Environ() {
			fmt.Println(pair)
		}
	}

	_, localRun := os.LookupEnv("LOCAL_RUN")
	if localRun {
		log.Println("Running on local environment")
		err := handler(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	// This is the normal entry point for lambdas
	lambda.Start(handler)
}
