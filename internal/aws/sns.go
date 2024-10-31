package aws

import (
	"context"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

func LoadSnsConfig(ctx context.Context, region string) (*sns.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		slog.Error("failed to load AWS sdk config", "error", err)
		return nil, err
	}
	snsClient := sns.NewFromConfig(cfg)
	return snsClient, nil

}

func CreateSnsTopicARN(ctx context.Context, topicName string, client *sns.Client) (string, error) {
	createdTopicOutput, err := client.CreateTopic(ctx, &sns.CreateTopicInput{
		Name: aws.String(topicName),
	})
	if err != nil {
		slog.Error("failed to create topic: ", "error", err)
	}
	topicArn := *createdTopicOutput.TopicArn

	return topicArn, nil

}
