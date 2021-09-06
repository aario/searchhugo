package hugohelpers

import (
    "strings"
    "net/url"
    "path/filepath"
)

// AbsURL creates an absolute URL from the relative path given and the BaseURL set in config.
func AbsURL(path string, contentPath string, baseUrl string) string {
    dir := strings.Replace(
        filepath.Dir(path),
        contentPath,
        "",
        1,
    )

    if (strings.HasPrefix(dir, "/")) {
        dir = dir[1:]
    }

    baseName := strings.Replace(
        filepath.Base(path),
        ".md",
        "/",
        -1,
    )
    if (baseName == "_index/") {
        baseName = ""
    }

    return baseUrl + strings.Replace(
        strings.Replace(
            url.QueryEscape(
                strings.ToLower(
                    strings.Replace(
                        dir + "/" + baseName,
                        " ",
                        "-",
                        -1,
                    ),
                ),
            ),
            "%2F",
            "/",
            -1,
        ),
        "//",
        "/",
        -1,
    )
}
