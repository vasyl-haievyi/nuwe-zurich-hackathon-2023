package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, s3Event events.S3Event) error {
	sdkConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Printf("failed to load default config: %s\n", err)
		return fmt.Errorf("failed to load default config: %w", err)
	}
	s3Client := s3.NewFromConfig(sdkConfig)
	dynamoClient := dynamodb.NewFromConfig(sdkConfig)

	for _, record := range s3Event.Records {
		bucket := record.S3.Bucket.Name
		key := record.S3.Object.URLDecodedKey
		getResponse, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: &bucket,
			Key:    &key,
		})
		if err != nil {
			log.Printf("error getting an object %s/%s: %s\n", bucket, key, err.Error())
			return fmt.Errorf("error getting an object %s/%s: %w", bucket, key, err)
		}
		defer getResponse.Body.Close()
		object, err := io.ReadAll(getResponse.Body)
		if err != nil {
			return fmt.Errorf("could not read response body: %w", err)
		}
		log.Printf("received and object %q\n", object)

		clients := []Client{}
		if err := json.Unmarshal(object, &clients); err != nil {
			return fmt.Errorf("could not unmarchall client list: %w", err)
		}

		tableName, ok := os.LookupEnv("DATA_TABLE_NAME")
		if !ok || tableName == "" {
			return fmt.Errorf("env var %q not set or empty", "DATA_TABLE_NAME")
		}

		for _, client := range clients {
			item, err := attributevalue.MarshalMap(client)
			if err != nil {
				return fmt.Errorf("could not marshal client to item: %w", err)
			}

			_, err = dynamoClient.PutItem(ctx, &dynamodb.PutItemInput{
				TableName: &tableName,
				Item:      item,
			})

			if err != nil {
				return fmt.Errorf("could not put item: %w", err)
			}
		}
	}

	return nil
}

type Client struct {
	ID        string    `json:"id" dynamodbav:"id"`
	Name      string    `json:"name" dynamodbav:"name"`
	Surname   string    `json:"surname" dynamodbav:"surname"`
	Birthdate string    `json:"birthdate" dynamodbav:"birthdate"`
	Address   string    `json:"address" dynamodbav:"address"`
	Car       ClientCar `json:"car" dynamodbav:"car"`
	Fee       int       `json:"fee" dynamodbav:"fee"`
}

type ClientCar struct {
	Make         string `json:"make" dynamodbav:"make"`
	Model        string `json:"model" dynamodbav:"model"`
	Year         int    `json:"year" dynamodbav:"year"`
	Color        string `json:"color" dynamodbav:"color"`
	Plate        string `json:"plate" dynamodbav:"plate"`
	Mileage      int    `json:"mileage" dynamodbav:"mileage"`
	FuelType     string `json:"fuelType" dynamodbav:"fuelType"`
	Transmission string `json:"transmission" dynamodbav:"transmission"`
}
