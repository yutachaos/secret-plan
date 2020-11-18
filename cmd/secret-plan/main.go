package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/urfave/cli/v2"
	"github.com/yutachaos/secret-plan/internal/secret"
)

var version = "master"

func main() {
	app := NewApp()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// NewApp is the factory method to return secret-plan.
func NewApp() *cli.App {
	app := cli.NewApp()
	app.Name = "secret-plan"
	app.Usage = "For aws secretsmanager save,update,diff tool"
	app.Version = version
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
		if err != nil {
			return err
		}

		secretValue = string(b)
	}

	versionID := ctx.String("version-id")

	sec := secret.NewSecret()

	currentSecret, secretExist, err := sec.Get(secretName, versionID)
	if err != nil {
		return err
	}

	if !diff(secretName, currentSecret, secretValue) {
		fmt.Println("No change.")

		return nil
	}

	if approve() {
		err := sec.Save(secretName, secretValue, secretExist)
		if err != nil {
			return err
		}
	} else {
		fmt.Print("No Updated.")
	}

	return nil
}

func diff(secretName string, currentSecretValue string, secretValue string) (diff bool) {
	if currentSecretValue == secretValue {
		return false
	}

	dmp := diffmatchpatch.New()

	fmt.Printf("name: %s \n", secretName)

	diffs := dmp.DiffMain(currentSecretValue, secretValue, true)

	fmt.Println("------------------------------------------------------------------------")
	fmt.Println(dmp.DiffPrettyText(diffs))
	fmt.Println("------------------------------------------------------------------------")

	return true
}

func approve() bool {
	fmt.Println("Apply? (yes/no)")

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		if scanner.Text() != "yes" {
			return false
		}
	}

	return true
}
