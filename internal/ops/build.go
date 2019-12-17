package main

import (
	"github.com/spf13/cobra"
)

var cmdBuild = &cobra.Command{
	Use:   "build <stack> <env>",
	Short: "Controls the build process for the application",
	Args:  cobra.MinimumNArgs(2),
	Run:   runBuildCmd,
}

func runBuildCmd(cmd *cobra.Command, args []string) {

}

func init() {
	rootCmd.AddCommand(cmdBuild)
}

// ## build: run tests and compile full app in production mode
// build: FRONTEND_API_HOST="https://${APP_NAME}.clintjedwards.com"
// build: check-semver-included
// 	protoc --go_out=plugins=grpc:. api/*.proto
// 	protoc --js_out=import_style=commonjs,binary:./frontend/src/ --grpc-web_out=import_style=typescript,mode=grpcwebtext:./frontend/src/ -I ./api/ api/*.proto
// 	go mod tidy
// 	go test ./utils
// 	npm run --prefix ./frontend build:production
// 	packr build -ldflags $(GO_LDFLAGS) -o $(BUILD_PATH)

// ## build-backend: build backend without frontend assets
// build-backend: SEMVER=v0.0.1
// build-backend:
// 	protoc --go_out=plugins=grpc:. api/*.proto
// 	go mod tidy
// 	go test ./utils
// 	go build -ldflags $(GO_LDFLAGS) -o $(BUILD_PATH)

// ## build-dev: build development version of app
// build-dev: SEMVER=v0.0.1
// build-dev:
// 	npx webpack --config="./frontend/webpack.config.js" --mode="development"
// 	packr build -ldflags $(GO_LDFLAGS) -o $(BUILD_PATH)

// ## build-protos: build required protobuf files
// build-protos:
// 	protoc --go_out=plugins=grpc:. api/*.proto
// 	protoc --js_out=import_style=commonjs,binary:./frontend/src/ --grpc-web_out=import_style=typescript,mode=grpcwebtext:./frontend/src/ -I ./api/ api/*.proto

// func runUsersCreateCmd() {
// 	name := args[0]

// 	fmt.Printf("Password: ")
// 	pass, err := gopass.GetPasswdMasked()
// 	if err != nil {
// 		log.Fatalf("failed retrieve password")
// 	}

// 	hash, err := utils.HashPassword(pass)
// 	if err != nil {
// 		log.Fatalf("failed to hash password: %v", err)
// 	}

// 	storage, err := storage.InitStorage()
// 	if err != nil {
// 		log.Fatalf("could not connect to storage: %v", err)
// 	}

// 	err = storage.CreateUser(name, &api.User{
// 		Name: name,
// 		Hash: string(hash),
// 	})
// 	if err != nil {
// 		log.Fatalf("could not create user: %v", err)
// 	}
// }

// func init() {
// 	cmdUsers.AddCommand(cmdUsersCreate)
// }
