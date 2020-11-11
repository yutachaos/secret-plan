package secret

type Secret interface {
	Plan(name string, content string, versionId string) (secretExist bool, err error)
	Save(name string, content string, versionId string, secretExist bool) (err error)
}

func NewSecret() Secret {
	secret := newAws()
	return secret
}
