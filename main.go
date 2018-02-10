package main

import (
	"context"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2/clientcredentials"
	"slurpCwsImages/service"
	"slurpCwsImages/slurp"
)

func main() {
	client := buildOAuthClient()
	cwsService := service.BuildCwsService(client)
	awsService := buildS3Client()

	slurp.SlurpImages(cwsService, awsService)
}

func buildS3Client() service.AwsService {
	var bucket string
	var timeout time.Duration

	bucket = getEnv("BUCKET")
	timeout = time.Second * time.Duration(30)

	return service.BuildAwsService(bucket, timeout)
}

func buildOAuthClient() *http.Client {
	return buildOAuthClientConfig().Client(context.Background())
}

func buildOAuthClientConfig() *clientcredentials.Config {
	return &clientcredentials.Config{
		ClientID:     getEnv("CWS_CLIENT_ID"),
		ClientSecret: getEnv("CWS_CLIENT_SECRET"),
		TokenURL:     getEnv("CWS_CLIENT_TOKEN"),
		Scopes:       []string{},
	}
}

func getEnv(key string) string {
	value, isPresent := os.LookupEnv(key)
	if !isPresent {
		log.WithField("key", key).Error("Env variable not set")
		panic(1)
	}

	return value
}
