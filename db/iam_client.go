package db

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
	"github.com/shashank-mugiwara/laughingtale/conf"
	"go.uber.org/zap"
)

type IAM struct {
	laughingtaleLogger *zap.SugaredLogger
}

func NewIAMClient() *IAM {
	return &IAM{}
}

func (iamHandler *IAM) GetIamRdsCredential(ctx context.Context, host string) (string, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		iamHandler.laughingtaleLogger.Info("Could not create config for pgxpoll. Exiting ...")
		return "", err
	}

	// 840000 milliseconds = 14 minutes
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 840000*time.Millisecond)
	defer cancel()

	var dbUser string = conf.PostgresSetting.Username
	var dbHost string = host
	var dbPort string = conf.PostgresSetting.Port
	var dbEndpoint string = fmt.Sprintf("%s:%s", dbHost, dbPort)
	var region string = conf.AWSSetting.Region

	authenticationToken, err := auth.BuildAuthToken(ctxWithTimeout, dbEndpoint, region, dbUser, cfg.Credentials)
	if err != nil {
		iamHandler.laughingtaleLogger.Info("Could not retrieve authenticationToken from AWS. Exiting ...")
		return "", err
	}

	return authenticationToken, nil
}
