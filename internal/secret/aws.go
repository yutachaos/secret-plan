package secret

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
)

type aws struct {
	client       secretsmanageriface.SecretsManagerAPI
	versionStage string
}

func newAws(ctx *cli.Context) *aws {
	opts := session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}

	sess, err := session.NewSessionWithOptions(opts)
	if err != nil {
		panic(err)
	}

	client := secretsmanager.New(sess)
	versionStage := ctx.String("version-stage")

	return &aws{
		client:       client,
		versionStage: versionStage,
	}
}

func (a *aws) Get(name string, versionID string) (currentSecret string, err error) {
	input := secretsmanager.GetSecretValueInput{
		SecretId: &name,
	}
	if versionID != "" {
		input.VersionId = &versionID
	}

	if a.versionStage != "" {
		input.VersionStage = &a.versionStage
	}

	_, err = a.client.DescribeSecret(&secretsmanager.DescribeSecretInput{
		SecretId: &name,
	})

	if err != nil {
		return "", err
	}

	var currentSecretValue string

	secretValueOutput, err := a.client.GetSecretValue(&input)
	if err != nil {
		switch err.(type) {
		case *secretsmanager.ResourceNotFoundException:
			break
		default:
			return "", err
		}
	} else {
		currentSecretValue = *secretValueOutput.SecretString
	}

	return currentSecretValue, nil
}

func (a *aws) Save(name string, content string) (err error) {
	input := secretsmanager.PutSecretValueInput{SecretId: &name, SecretString: &content}
	if a.versionStage != "" {
		input.VersionStages = []*string{&a.versionStage}
	}

	output, err := a.client.PutSecretValue(&input)
	if err != nil {
		return err
	}

	versionID := *output.VersionId
	fmt.Printf("Update. Version: %s \n", versionID)

	return nil
}
