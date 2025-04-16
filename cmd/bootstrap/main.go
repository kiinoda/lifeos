package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
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
