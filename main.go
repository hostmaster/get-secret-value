package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

func main() {
	var id, stage, region string

	flag.StringVar(&id, "secret-id", "", "Secret Id")
	flag.StringVar(&stage, "version-stage", "AWSCURRENT", "Specifies the secret version that you want to retrieve (default: AWSCURRENT)")
	flag.StringVar(&region, "region", "us-east-1", "AWS Region (default: us-east-1)")
	flag.Parse()

	if len(id) == 0 {
		fmt.Println("Error: No secret id")
		os.Exit(1)
	}

	s, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String(region),
		},
	})

	sm := secretsmanager.New(s)
	output, err := sm.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(id),
		VersionStage: aws.String(stage),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeResourceNotFoundException:
				fmt.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			case secretsmanager.ErrCodeInvalidParameterException:
				fmt.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())
			case secretsmanager.ErrCodeInvalidRequestException:
				fmt.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
			case secretsmanager.ErrCodeDecryptionFailure:
				fmt.Println(secretsmanager.ErrCodeDecryptionFailure, aerr.Error())
			case secretsmanager.ErrCodeInternalServiceError:
				fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(*output.SecretString)
}
