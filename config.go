package main

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
)

type Config struct {
    Username string `json:"username"`
    Cookie   string `json:"cookie"`
}

func LoadConfig(filePath string) (*Config, error) {
    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        // Create the config directory if it doesn't exist
        configDir := filepath.Dir(filePath)
        err = os.MkdirAll(configDir, os.ModePerm)
        if err != nil {
            return nil, fmt.Errorf("error creating config directory: %v", err)
        }

        // Create a blank config file
        blankConfig := &Config{
            Username: "",
            Cookie:   "",
        }
        file, err := os.Create(filePath)
        if err != nil {
            return nil, fmt.Errorf("error creating config file: %v", err)
        }
        defer file.Close()

        encoder := json.NewEncoder(file)
        err = encoder.Encode(blankConfig)
        if err != nil {
            return nil, fmt.Errorf("error encoding blank config: %v", err)
        }

        fmt.Printf("Config file not found. A blank config has been created at %s\n", filePath)
        fmt.Println("Please edit the config file and provide the necessary values.")
        os.Exit(1)
    }

    file, err := os.Open(filePath)
    if err != nil {
        return nil, fmt.Errorf("error opening config file: %v", err)
    }
    defer file.Close()

    var config Config
    decoder := json.NewDecoder(file)
    err = decoder.Decode(&config)
    if err != nil {
        return nil, fmt.Errorf("error decoding config file: %v", err)
    }

    if config.Username == "" || config.Cookie == "" {
        return nil, fmt.Errorf("missing required config values")
    }

    return &config, nil
}
