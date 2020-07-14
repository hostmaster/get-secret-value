package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

func main() {
	var id, stage, region string

	flag.StringVar(&id, "secret-id", "", "Secret Id")
	flag.StringVar(&stage, "version-stage", "AWSCURRENT", "the staging label attached to the version.(default: AWSCURRENT)")
	flag.StringVar(&region, "region", "", "AWS Region")
	flag.Parse()

	if len(id) == 0 {
		fmt.Println("Error: No secret ID specified")
		os.Exit(1)
	}

	s, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String(region),
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	sm := secretsmanager.New(s)
	output, err := sm.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(id),
		VersionStage: aws.String(stage),
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(*output.SecretString)
}
