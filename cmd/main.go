package main

import "github.com/ArtemChadaev/GoCreateHistory/internal/config"

func main() {
	_, err := config.Load()

	if err != nil {

	}
}
