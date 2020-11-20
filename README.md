# secret-plan

![CI](https://github.com/yutachaos/secret-plan/workflows/CI/badge.svg)

For secret value diff and save tool.(e.g aws secretsmanager)

## Usage. 
- --secret-name
   - Secret's name 
- --secret-value 
   - Value to be registered in Secret
- --version-id
   - Specify the version-id to be acquired (optional).
- --is-file
   - When you specify a string or a file, or enable this flag, the file is recognized as a file path and read.

## Demo. 

![demo](./demo.gif)

## Requirement (Current AWS support only)
- The configuration from environment variables and requires an AWS authentication key

### Run
- go run ./cmd/secret-plan/main.go --secret-name test/diff --secret-value "update secret"

### Build
- make build

### test 
- make test 
