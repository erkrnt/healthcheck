package main

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/joho/godotenv"
)

// LoadDotEnv : See docs https://github.com/joho/godotenv#precendence--conventions
func LoadDotEnv() {
	env := os.Getenv("HEALTHCHECK_ENV")
	if "" == env {
		env = "development"
	}
	godotenv.Load(".env." + env + ".local")
	if "test" != env {
		godotenv.Load(".env.local")
	}
	godotenv.Load(".env." + env)
	godotenv.Load()
}

func main() {
	LoadDotEnv()

	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, c := range containers {
		fmt.Printf("%s %s\n", c.ID[:10], c.Image)
	}
}
