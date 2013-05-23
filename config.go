package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	User  string
	Token string
}

func LoadConfig(filename string) (*Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return &Config{}, err
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	dec := json.NewDecoder(reader)

	var c Config
	err = dec.Decode(&c)
	if err != nil {
		log.Fatal(err)
	}

	return &c, nil
}
