package config

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	awsssm "github.com/PaddleHQ/go-aws-ssm"
	"github.com/aws/aws-sdk-go/aws"
)

type ctxKey int

type AppConfig struct {
	SpreadsheetId     string          `json:"spreadsheetId"`
	ServiceAccountKey json.RawMessage `json:"serviceAccountKey"`
	Sender            string          `json:"sender"`
	Recipient         string          `json:"recipient"`
}

func (c *AppConfig) ServiceAccountKeyBytes() ([]byte, error) {
	if len(c.ServiceAccountKey) == 0 {
		return nil, fmt.Errorf("serviceAccountKey not found")
	}
	return []byte(c.ServiceAccountKey), nil
}

const configKey ctxKey = 0

func NewConfig(region string, fullSSMPath string) (*AppConfig, error) {
	pathComponents := strings.Split(fullSSMPath, "/")
	particle := pathComponents[len(pathComponents)-1]
	path := strings.TrimSuffix(fullSSMPath, particle)
	cfg := &AppConfig{}
	awsConfig := &aws.Config{
		Region: aws.String(region),
	}
	pmstore, err := awsssm.NewParameterStore(awsConfig)
	if err != nil {
		return cfg, err
	}

	//Requesting the base path
	params, err := pmstore.GetAllParametersByPath(path, true)
	if err != nil {
		return cfg, err
	}

	//And getting a specific value
	value := params.GetValueByName(particle)
	err = json.Unmarshal([]byte(value), cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}

func ConfigFromContext(ctx context.Context) (*AppConfig, error) {
	if cfg, ok := ctx.Value(configKey).(*AppConfig); ok {
		return cfg, nil
	}

	return &AppConfig{}, errors.New("Error retrieving config from context")
}

func ContextWithConfig(ctx context.Context, cfg *AppConfig) context.Context {
	return context.WithValue(ctx, configKey, cfg)
}
