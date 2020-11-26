package secret

import "github.com/urfave/cli/v2"

type Secret interface {
	Get(name string, versionID string) (currentSecret string, err error)
	Save(name string, content string) (err error)
}

func NewSecret(ctx *cli.Context) Secret {
	secret := newAws(ctx)

	return secret
}
