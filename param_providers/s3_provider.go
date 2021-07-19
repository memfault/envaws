package param_providers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type s3Provider struct {
	bucket       string
	objectKey    string
	wantedParams []string
	client       *s3.Client
	hash         string
}

func (s3Provider s3Provider) getAndHashData() string {
	result, err := s3Provider.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s3Provider.bucket),
	})
	if err != nil {
		log.Fatal(err)
	}
	defer result.Body.Close()
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(body)
	var jsonResult map[string]interface{}
	err = json.Unmarshal([]byte(bodyString), &jsonResult)
	if err != nil {
		log.Fatal(err)
	}
	wantedParams := ForceMapValuesToString(FilterParams(jsonResult, s3Provider.wantedParams))
	return HashParams(wantedParams)
}

func (s3Provider *s3Provider) Init() {
	initialHash := s3Provider.getAndHashData()
	s3Provider.hash = initialHash
}

func (s3Provider s3Provider) Changed() bool {
	currentHash := s3Provider.getAndHashData()
	return s3Provider.hash != currentHash
}

func NewS3Provider() *s3Provider {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	s3client := s3.NewFromConfig(cfg)

	return &s3Provider{
		bucket:       "sanko",
		objectKey:    "foo/bar.json",
		wantedParams: []string{"wow"},
		client:       s3client,
	}
}
