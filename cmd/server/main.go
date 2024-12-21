package main

import (
	"fmt"

	"github.com/mateusfaustino/go-rest-api-III/configs"
)

func main() {
	config, _:= configs.LoadConfig(".")
	fmt.Println(config.DBDriver)
}