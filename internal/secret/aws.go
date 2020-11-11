package secret

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/sergi/go-diff/diffmatchpatch"
)

type aws struct {
	client *secretsmanager.SecretsManager
}

func (a aws) Plan(name string, content string, versionId string) (secretExist bool, err error) {
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
			return secretExist, err
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
				return secretExist, err
			}
		} else {
			currentSecretValue = *secretValueOutput.SecretString
		}
	}

	if currentSecretValue == content {
		fmt.Println("No change.")
		return secretExist, nil
	}
	dmp := diffmatchpatch.New()

	fmt.Printf("name: %s \n", name)
	diffs := dmp.DiffMain(currentSecretValue, content, true)
	fmt.Println("------------------------------------------------------------------------")
	fmt.Println(dmp.DiffPrettyText(diffs))
	fmt.Println("------------------------------------------------------------------------")
	return secretExist, nil
}

func (a aws) Save(name string, content string, versionId string, secretExist bool) (err error) {
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

func newAws() aws {
	opts := session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}

	sess, err := session.NewSessionWithOptions(opts)

	if err != nil {
		panic(err)
	}
	a := secretsmanager.New(sess)
	return aws{
		client: a,
	}

}
