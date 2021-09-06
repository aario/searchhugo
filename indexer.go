package main

import (
    "io/ioutil"
    "strings"
    "regexp"
)

const (
    MIN_WORD_LENGTH = 1
    CONTENT_SUMMARY_LENGTH = 300
    SPLIT_REGEXP_DELIMITER_PATTERN = "[^a-zA-Z0-9\\-]"
)

type Post struct {
    Url string
    Title string
    Date string
    Summary string
}

type Occurance struct {
    Score int
    Position int
}

type Index struct {
    //Words is map[word][postId][scores...]
    Words map[string]map[int][]Occurance
    Posts []Post
}

var (
    index Index
    splitRegexp *   regexp.Regexp
)

func initializeIndex() {
    index.Posts = make([]Post, 0)
    index.Words = make(map[string]map[int][]Occurance, 0)
    splitRegexp = regexp.MustCompile(SPLIT_REGEXP_DELIMITER_PATTERN)
}

func stripNonAlphaNumeric(s string) string {
    var result strings.Builder
    for i := 0; i < len(s); i++ {
        b := s[i]
        if ('a' <= b && b <= 'z') ||
            ('A' <= b && b <= 'Z') ||
            ('0' <= b && b <= '9') ||
            b == ' ' {
            result.WriteByte(b)
        }
    }
    return result.String()
}

func shouldSkipWord(word string) bool {
    if (len(word) < MIN_WORD_LENGTH) {
        return true
    }

    _, mustSkip := STOP_WORDS[word]
    if (mustSkip) {
        return true
    }

    return false
}

func headerLookup(headerLines []string, name string) string {
    for _, headerLine := range headerLines {
        if (!strings.HasPrefix(headerLine, name + " = ")) {
            continue
        }

        headerLineParts := strings.Split(headerLine, "\"")
        if (len(headerLineParts) < 3) {
            panic("Line in header does not contain string quoted with '\"': " + headerLine)
        }

        return headerLineParts[1]
    }

    return ""
}

func readContent(path string, url string) (Post, string) {
    dat, err := ioutil.ReadFile(path)
    checkError(err)

    parts := strings.Split(string(dat), "+++")
    if (len(parts) < 3) {
        panic("File " + path + " does not contain header section start and end with +++")
    }

    header := parts[1]
    headerLines := strings.Split(header, "\n")
    title := headerLookup(headerLines, "title")
    date := headerLookup(headerLines, "date")
    content := strings.Trim(
        strings.Trim(
            parts[2],
            " ",
        ),
        "\n",
    )

    summary := content
    if (len(content) > CONTENT_SUMMARY_LENGTH) {
        summary = content[:CONTENT_SUMMARY_LENGTH]
    }

    post := Post{
        Url: url,
        Title: title,
        Date: date,
        Summary: summary,
    }


    return post, content
}

func addWordOccuranceToIndex(
    postId int,
    word string,
    positionScore int,
    positionIndex int,
) {
    word = stripNonAlphaNumeric(
        word,
    )
    if (shouldSkipWord(word)) {
        return
    }

    wordOccurances, ok := index.Words[word]
    if (!ok) {
        wordOccurances = make(map[int][]Occurance, 0)
    }

    occurances, ok := index.Words[word][postId]
    if (!ok) {
        occurances = make([]Occurance, 0)
    }

    occurance := Occurance {
        Score : positionScore,
        Position: positionIndex,
    }
    occurances = append(occurances, occurance)
    wordOccurances[postId] = occurances
    index.Words[word] = wordOccurances
}

func indexWord(
    postId int,
    word string,
    positionScore int,
    positionIndex int,
) {
    addWordOccuranceToIndex(
        postId,
        word,
        positionScore * RATING_SCORE_SAME,
        positionIndex,
    )

    synonyms, antonyms := thesaurusLookup(word)

    for _, synonym := range synonyms {
        addWordOccuranceToIndex(
            postId,
            synonym,
            positionScore * RATING_SCORE_SYNONYM,
            positionIndex,
        )
    }

    for _, antonym := range antonyms {
        addWordOccuranceToIndex(
            postId,
            antonym,
            positionScore * RATING_SCORE_ANTONYM,
            positionIndex,
        )
    }
}

/*
 * returns int the last positionIndex counted
 */
func indexText(
    postId int,
    text string,
    positionScore int,
    beginPositionIndex int,
) int {
    positionIndex := beginPositionIndex
    for _, word := range splitRegexp.Split(text, -1) {
        indexWord(
            postId,
            strings.ToLower(word),
            positionScore,
            positionIndex,
        )
        positionIndex++
    }

    return positionIndex
}

func indexFile(path string, url string) {
    postId := len(index.Posts)

    post, content := readContent(path, url)
    index.Posts = append(index.Posts, post)

    contentBeginPositionIndex := indexText(
        postId,
        post.Title,
        RATING_SCORE_TITLE,
        0,
    )
    indexText(
        postId,
        content,
        RATING_SCORE_CONTENT,
        contentBeginPositionIndex,
    )
}
