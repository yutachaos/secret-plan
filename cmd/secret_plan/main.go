package main

import (
	"bufio"
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/yutachaos/secret-plan/internal/secret"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	app := NewApp()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// NewApp is the factory method to return secret-plan
func NewApp() *cli.App {

	app := cli.NewApp()
	app.Name = "secret-plan"
	app.Usage = "For aws secretsmanager save,update,diff tool"
	app.Version = "0.0.1"
	app.ArgsUsage = "target"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     "secret-name",
			Usage:    "Specifies secret name",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "secret-value",
			Usage:    "Specifies text data that you want to encrypt and store in this new version of the secret",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "version-id",
			Usage:    "Whether to specify version-id",
			Required: false,
		},
		&cli.BoolFlag{
			Name:     "is-file",
			Usage:    "Specify filepath or string",
			Required: false,
		},
	}
	app.Action = Run
	return app
}

// Run is execute secret-plan.
func Run(ctx *cli.Context) (err error) {
	secretName := ctx.String("secret-name")
	secretValue := ctx.String("secret-value")
	isFile := ctx.Bool("is-file")

	if isFile {
		f, err := os.Open(secretValue)
		if err != nil {
			return err
		}
		defer f.Close()

		b, err := ioutil.ReadAll(f)
		secretValue = string(b)
	}
	versionId := ctx.String("version-id")

	secret := secret.NewSecret()

	secretExist, err := secret.Plan(secretName, secretValue, versionId)
	if err != nil {
		return err
	}

	approved := approve()

	if approved {
		err := secret.Save(secretName, secretValue, versionId, secretExist)
		if err != nil {
			return err
		}
	} else {
		fmt.Print("No Updated.")
	}
	return nil
}

func approve() bool {
	fmt.Println("Apply? (yes/no)")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if scanner.Text() == "yes" {
			return true
		} else {
			break
		}
	}
	return false
}
