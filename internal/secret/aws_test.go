package secret

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
	"github.com/stretchr/testify/assert"
	"tideland.dev/go/audit/capture"
)

type fakeAws struct {
	secretsmanageriface.SecretsManagerAPI
	fakeDescribeSecret func(*secretsmanager.DescribeSecretOutput, error) (*secretsmanager.DescribeSecretOutput, error)
	fakeGetSecretValue func(*secretsmanager.GetSecretValueOutput, error) (*secretsmanager.GetSecretValueOutput, error)
	fakePutSecretValue func(*secretsmanager.PutSecretValueOutput, error) (*secretsmanager.PutSecretValueOutput, error)
}

type describeSecretOutput struct {
	output *secretsmanager.DescribeSecretOutput
	err    error
}

type getSecretValueOutput struct {
	output *secretsmanager.GetSecretValueOutput
	err    error
}

type putSecretValueOutput struct {
	output *secretsmanager.PutSecretValueOutput
	err    error
}

func (f fakeAws) DescribeSecret(*secretsmanager.DescribeSecretInput) (output *secretsmanager.DescribeSecretOutput, err error) {
	return f.fakeDescribeSecret(output, err)
}

func (f fakeAws) GetSecretValue(*secretsmanager.GetSecretValueInput) (output *secretsmanager.GetSecretValueOutput, err error) {
	return f.fakeGetSecretValue(output, err)
}

func (f fakeAws) PutSecretValue(*secretsmanager.PutSecretValueInput) (output *secretsmanager.PutSecretValueOutput, err error) {
	return f.fakePutSecretValue(output, err)
}

func TestGet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		secretName         string
		versionID          string
		mockDescribeSecret describeSecretOutput
		mockGetSecretValue getSecretValueOutput
		err                error
		currentSecret      string
	}{
		{
			secretName: "name1",
			versionID:  "",
			mockDescribeSecret: describeSecretOutput{
				output: &secretsmanager.DescribeSecretOutput{},
				err:    nil,
			},
			mockGetSecretValue: getSecretValueOutput{
				output: &secretsmanager.GetSecretValueOutput{SecretString: toStrPtr("name2")},
				err:    nil,
			},
			err:           nil,
			currentSecret: "name2",
		},
		{
			secretName: "name2",
			versionID:  "",
			mockDescribeSecret: describeSecretOutput{
				output: &secretsmanager.DescribeSecretOutput{},
				err:    &secretsmanager.ResourceNotFoundException{},
			},
			mockGetSecretValue: getSecretValueOutput{
				output: &secretsmanager.GetSecretValueOutput{SecretString: toStrPtr("")},
				err:    nil,
			},
			err:           &secretsmanager.ResourceNotFoundException{},
			currentSecret: "",
		},
		{
			secretName: "name3",
			versionID:  "",
			mockDescribeSecret: describeSecretOutput{
				output: &secretsmanager.DescribeSecretOutput{},
				err:    nil,
			},
			mockGetSecretValue: getSecretValueOutput{
				output: &secretsmanager.GetSecretValueOutput{SecretString: toStrPtr("")},
				err:    &secretsmanager.ResourceNotFoundException{},
			},
			err:           nil,
			currentSecret: "",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(fmt.Sprintf("Test Get: %s", test.secretName), func(t *testing.T) {
			t.Parallel()
			mockAws := aws{
				client: fakeAws{
					fakeDescribeSecret: func(*secretsmanager.DescribeSecretOutput, error) (*secretsmanager.DescribeSecretOutput, error) {
						return test.mockDescribeSecret.output, test.mockDescribeSecret.err
					},
					fakeGetSecretValue: func(*secretsmanager.GetSecretValueOutput, error) (*secretsmanager.GetSecretValueOutput, error) {
						return test.mockGetSecretValue.output, test.mockGetSecretValue.err
					},
				},
			}
			currentSecret, err := mockAws.Get(test.secretName, test.versionID)

			assert.Equal(t, test.err, err)
			assert.Equal(t, test.currentSecret, currentSecret)
		})
	}
}

func TestSave(t *testing.T) {
	t.Parallel()

	tests := []struct {
		secretName         string
		secretValue        string
		mockPutSecretValue putSecretValueOutput
		err                error
		stdout             string
	}{
		{
			secretName: "name1",
			mockPutSecretValue: putSecretValueOutput{
				output: &secretsmanager.PutSecretValueOutput{VersionId: toStrPtr("versionID1")},
				err:    nil,
			},
			err:    nil,
			stdout: "Update. Version: versionID1 \n",
		},
		{
			secretName: "name2",
			mockPutSecretValue: putSecretValueOutput{
				output: &secretsmanager.PutSecretValueOutput{VersionId: toStrPtr("versionID2")},
				err:    nil,
			},
			err:    nil,
			stdout: "Update. Version: versionID2 \n",
		},
		{
			secretName: "name3",
			mockPutSecretValue: putSecretValueOutput{
				output: &secretsmanager.PutSecretValueOutput{VersionId: toStrPtr("versionID3")},
				err:    &secretsmanager.InvalidRequestException{},
			},
			err: &secretsmanager.InvalidRequestException{},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(fmt.Sprintf("Test Save: %s", test.secretName), func(t *testing.T) {
			t.Parallel()
			mockAws := aws{
				client: fakeAws{
					fakePutSecretValue: func(*secretsmanager.PutSecretValueOutput, error) (*secretsmanager.PutSecretValueOutput, error) {
						return test.mockPutSecretValue.output, test.mockPutSecretValue.err
					},
				},
			}
			var err error
			stdout := capture.Stdout(func() {
				err = mockAws.Save(test.secretName, test.secretValue)
			})
			assert.Equal(t, test.err, err)
			assert.Equal(t, test.stdout, stdout.String())
		})
	}
}

func toStrPtr(s string) *string {
	return &s
}
