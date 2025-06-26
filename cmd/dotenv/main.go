package main

import (
	"log"

	"github.com/fatih/color"
	"github.com/taimats/internal"
)

func main() {
	err := internal.LoadEnvFile("./cmd/dotenv/test.env")
	if err != nil {
		log.Fatal(err)
	}
	color.Green("環境変数の設定に成功!!")
}
