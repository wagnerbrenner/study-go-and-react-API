package main

import (
	"fmt"
	"os/exec"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(fmt.Sprintf("Error loading .env file: %v", err))
	}

	cmd := exec.Command(
		"tern",
		"migrate",
		"--migrations",
		"./internal/store/pgstore/migrations",
		"--config",
		"./internal/store/pgstore/migrations/tern.conf",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		panic(fmt.Sprintf("Error running tern command: %v\nOutput: %s", err, string(output)))
	}
	fmt.Printf("Command output: %s\n", string(output))
}
