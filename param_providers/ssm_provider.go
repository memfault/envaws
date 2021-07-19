package param_providers

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type ssmProvider struct {
	hash         string
	wantedParams []string
	client       *ssm.Client
}

func (ssmProvider ssmProvider) getAndHashData() string {
	result, err := ssmProvider.client.GetParameters(context.TODO(), &ssm.GetParametersInput{
		Names:          ssmProvider.wantedParams,
		WithDecryption: false,
	})
	if err != nil {
		log.Fatal(err)
	}

	params := make(map[string]string)
	for _, p := range result.Parameters {
		params[*p.Name] = *p.Value
	}
	return HashParams(params)
}

func (ssmProvider *ssmProvider) Init() {
	initialHash := ssmProvider.getAndHashData()
	ssmProvider.hash = initialHash
}

func (ssmProvider ssmProvider) Changed() bool {
	currentHash := ssmProvider.getAndHashData()
	return ssmProvider.hash != currentHash
}

func NewSSMProvider(wantedParams []string) *ssmProvider {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	return &ssmProvider{
		wantedParams: wantedParams,
		client:       ssm.NewFromConfig(cfg),
	}
}
