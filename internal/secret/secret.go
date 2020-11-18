package secret

type Secret interface {
	Get(name string, versionID string) (currentSecret string, secretExist bool, err error)
	Save(name string, content string, secretExist bool) (err error)
}

func NewSecret() Secret {
	secret := newAws()

	return secret
}
