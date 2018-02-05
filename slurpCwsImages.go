package main

import (
	"net/http"
	"context"
	"os"
	"strings"
	"time"
	"fmt"

	"golang.org/x/oauth2/clientcredentials"
	log "github.com/sirupsen/logrus"
	"slurpCwsImages/service"
)

func main() {
	client := buildOAuthClient()
	cwsService := service.BuildCwsService(client)
	awsService := buildS3Client()

	resp, err := cwsService.GetProducts()
	if err != nil {
		panic(1)
	}

	for _, item := range resp.Items {
		for _, region := range item.Regions {
			if strings.ToUpper(region) == "WORLDWIDE" {
				for _, image := range item.Images {
					file, fileExt, err := cwsService.GetProductImage(image.Image)
					if err != nil {
						continue
					}

					awsService.UploadItemToS3(file, fileExt, fmt.Sprintf("%s_%s%s", item.ProductID, strings.ToLower(image.Format), fileExt))
				}
			}
		}
	}
}

func buildS3Client() service.AwsService {
	var bucket string
	var timeout time.Duration

	bucket = getEnv("BUCKET")
	timeout = time.Second * time.Duration(20)

	return service.BuildAwsService(bucket, timeout)
}

func buildOAuthClient() *http.Client {
	return buildOAuthClientConfig().Client(context.Background())
}

func buildOAuthClientConfig() *clientcredentials.Config {
	return &clientcredentials.Config{
		ClientID: getEnv("CWS_CLIENT_ID"),
		ClientSecret: getEnv("CWS_CLIENT_SECRET"),
		TokenURL: getEnv("CWS_CLIENT_TOKEN"),
		Scopes: []string{},
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



