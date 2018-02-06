package main

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2/clientcredentials"
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
					fileExt, err := cwsService.HeadProductImage(image.Image)
					if err != nil {
						continue
					}

					if !awsService.S3ItemExists(item.ProductID, fileExt, image.Format) {
						file, _, err := cwsService.GetProductImage(image.Image)
						if err != nil {
							continue
						}

						awsService.UploadItemToS3(file, fileExt, item.ProductID, image.Format)
					}
				}
			}
		}
	}
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
