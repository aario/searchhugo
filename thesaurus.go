package main

import (
    "fmt"
    "strings"
)

const (
    QUERY = "select synonyms, antonyms from thesaurus where word = \"%s\""
)

func splitWords(words string) []string {
    return strings.Split(words, ",")
}

func appendCommaSeparatedWords(words []string, commaSeparatedWords string) []string {
    for _, word := range splitWords(commaSeparatedWords) {
        word := strings.TrimSpace(word)
        if len(word) == 0 {
            continue
        }

        words = append(words, word)
    }

    return words
}

func thesaurusLookup(word string) ([]string, []string) {
    synonyms := make([]string, 0)
    antonyms := make([]string, 0)
    query := fmt.Sprintf(QUERY, word)
    fmt.Println("Query: " + query)
    rows := dbSelect(query)

    for _, row := range rows {
        synonyms = appendCommaSeparatedWords(synonyms, row[0])

        antonyms = appendCommaSeparatedWords(antonyms, row[1])
    }

    return synonyms, antonyms
}
