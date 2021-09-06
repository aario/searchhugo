package main

import (
    "flag"
    "fmt"
    "path/filepath"
    "os"
    "io/ioutil"
    "log"
    "github.com/aario/searchhugo/hugohelpers"
    "encoding/json"
    "strings"
)

func checkError(e error) {
    if (e != nil) {
        panic(e)
    }
}

var (
    rootPath string
    baseUrl string
    contentPath string
    outputPath string
    absoluteContentPath string
)

func walkPath(
    path string,
    info os.FileInfo,
    err error,
) error {
    if err != nil {
        return err
    }
    if (info.IsDir()) {
        return nil
    }
    url := hugohelpers.AbsURL(path, absoluteContentPath, baseUrl)
    log.Println(path, url)

    indexFile(path, url)
    return nil
}

func main() {
    flag.StringVar(&rootPath, "path", ".", "Path to hugo site root folder")
    flag.StringVar(&baseUrl, "baseUrl", ".", "Base url of the hugo site")
    flag.StringVar(&outputPath, "json", "./index.json", "Path to json output file")

    flag.Parse()

    dbConnect()

    var err error
    var currentPath string

    if (strings.HasPrefix(outputPath, ".")) {
        currentPath, err = os.Getwd()
        outputPath, err = filepath.Abs(currentPath + "/" + outputPath)
        checkError(err)
    }

    contentPath = rootPath + "/content"
    checkError(os.Chdir(contentPath))

    absoluteContentPath, err = os.Getwd()
    checkError(err)
    fmt.Println("content path:", absoluteContentPath)

    initializeIndex()

    checkError(filepath.Walk(absoluteContentPath, walkPath))

    var jsonBytes []byte
    jsonBytes, err = json.Marshal(index)
    checkError(err)

    checkError(ioutil.WriteFile(outputPath, jsonBytes, 0644))
}
