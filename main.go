package main

import (
	"golang/tutorial/todo/cmd"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	cmd.Execute()
}
