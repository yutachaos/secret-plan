package secret

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
)

type aws struct {
	client secretsmanageriface.SecretsManagerAPI
}

func newAws() *aws {
	opts := session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}

	sess, err := session.NewSessionWithOptions(opts)

	if err != nil {
		panic(err)
	}
	client := secretsmanager.New(sess)
	return &aws{
		client: client,
	}

}

func (a *aws) Get(name string, versionId string) (currentSecret string, secretExist bool, err error) {
	input := secretsmanager.GetSecretValueInput{
		SecretId: &name,
	}
	if versionId != "" {
		input.VersionId = &versionId
	}

	_, err = a.client.DescribeSecret(&secretsmanager.DescribeSecretInput{
		SecretId: &name,
	})

	secretExist = true
	if err != nil {
		switch err.(type) {
		case *secretsmanager.ResourceNotFoundException:
			secretExist = false
			break
		default:
			return "", secretExist, err
		}
	}

	var currentSecretValue string
	if secretExist {
		secretValueOutput, err := a.client.GetSecretValue(&input)
		if err != nil {
			switch err.(type) {
			case *secretsmanager.ResourceNotFoundException:
				break
			default:
				return "", secretExist, err
			}
		} else {
			currentSecretValue = *secretValueOutput.SecretString
		}
	}
	return currentSecretValue, secretExist, nil

}

func (a *aws) Save(name string, content string, versionId string, secretExist bool) (err error) {
	if secretExist {
		output, err := a.client.PutSecretValue(&secretsmanager.PutSecretValueInput{SecretId: &name, SecretString: &content})
		if err != nil {
			return err
		}
		versionId = *output.VersionId
		fmt.Printf("Update. Version: %s \n", versionId)
	} else {

		output, err := a.client.CreateSecret(&secretsmanager.CreateSecretInput{Name: &name, SecretString: &content})
		if err != nil {
			return err
		}
		versionId = *output.VersionId
		fmt.Printf("Create. Version: %s \n", versionId)
	}
	return nil
}
