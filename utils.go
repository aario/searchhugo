package main

import (
    "os"
    "path/filepath"
)

func fileExists(path string) bool {
    _, err := os.Stat(path)

    return (err == nil)
}

func getExecutablePath() string {
    path, _ := filepath.Abs(filepath.Dir(os.Args[0]))

    return path
}
