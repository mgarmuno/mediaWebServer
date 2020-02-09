package utils

import (
	"path/filepath"
	"regexp"
	"strings"
)

const (
	sanitizeBrackets        = `\[(.*?)\]`
	sanitizeName            = `[^A-Za-z0-9\s]`
	sanitizeLoadClose       = `^[\s\p{Zs}]+|[\s\p{Zs}]+$`
	sanitizeRedundantSpaces = `[\s\p{Zs}]{2,}`
)

func GetMovieNameSessionEpisodeByFileName(fileName string) (string, string, string) {
	var ext string = filepath.Ext(fileName)
	fileName = fileName[0 : len(fileName)-len(ext)]
	fileName = strings.ReplaceAll(fileName, "_", " ")
	re := regexp.MustCompile(sanitizeBrackets)
	var movieName string = string(re.ReplaceAll([]byte(fileName), []byte("")))
	re = regexp.MustCompile(sanitizeName)
	movieName = string(re.ReplaceAll([]byte(movieName), []byte("")))
	movieName, seasson := getAndReplaceSessionByMovieName(movieName)
	movieName, episode := getAndReplaceEpisodeByMovieName(movieName)
	movieName = replaceSignsForSapacesAndTrim(movieName)
	return movieName, seasson, episode
}

func replaceSignsForSapacesAndTrim(fileName string) string {
	var sings []string = []string{"_", ".", "-"}
	for _, sign := range sings {
		fileName = strings.ReplaceAll(fileName, sign, " ")
	}
	return replaceRedundantSpacesAndTrim(fileName)
}

func replaceRedundantSpacesAndTrim(fileName string) string {
	reTrim := regexp.MustCompile(sanitizeLoadClose)
	reSpaces := regexp.MustCompile(sanitizeRedundantSpaces)
	fileName = reTrim.ReplaceAllString(fileName, "")
	fileName = reSpaces.ReplaceAllString(fileName, " ")
	return fileName
}

func getAndReplaceSessionByMovieName(movieName string) (string, string) {
	var seassons []string = []string{`(?i)s\s*\d{2}?`, `(?i)session\s*\d{2}?`}
	for _, seasson := range seassons {
		re := regexp.MustCompile(seasson)
		var matched string = string(re.Find([]byte(movieName)))
		if matched != "" {
			movieName = re.ReplaceAllString(movieName, "")
			return movieName, matched[len(matched)-2 : len(matched)]
		}
	}
	return movieName, ""
}

func getAndReplaceEpisodeByMovieName(movieName string) (string, string) {
	var episodes []string = []string{`(?i)e\s*\d{2}?`, `(?i)episode\s*\d{2}?`}
	for _, episode := range episodes {
		re := regexp.MustCompile(episode)
		var matched string = string(re.Find([]byte(movieName)))
		if matched != "" {
			movieName = re.ReplaceAllString(movieName, "")
			return movieName, matched[len(matched)-2 : len(matched)]
		}
	}
	return movieName, ""
}
