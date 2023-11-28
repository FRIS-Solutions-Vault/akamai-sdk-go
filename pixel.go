package akamai

import (
	"errors"
	"io"
	"regexp"
	"strconv"
	"strings"
)

var (
	pixelBazaExpr              = regexp.MustCompile(`bazadebezolkohpepadr="(\d+)"`)
	pixelScriptExpr            = regexp.MustCompile(`src="(https?://.+/akam/\d+/\w+)"`)
	pixelScriptVarExpr         = regexp.MustCompile(`g=_\[(\d+)]`)
	pixelScriptStringArrayExpr = regexp.MustCompile(`var _=\[(.+)];`)
	pixelScriptStringsExpr     = regexp.MustCompile(`("[^",]*")`)

	ErrPixelBazaVarNotFound   = errors.New("akamai-sdk-go: pixel Baza var not found")
	ErrPixelScriptNotFound    = errors.New("akamai-sdk-go: script not found")
	ErrPixelScriptVarNotFound = errors.New("akamai-sdk-go: script var not found")
)

// ParsePixelBazaVar gets the required pixel challenge variable "bazadebezolkohpepadr" from the given HTML code src.
func ParsePixelBazaVar(reader io.Reader) (int, error) {
	src, err := io.ReadAll(reader)
	if err != nil {
		return 0, errors.Join(ErrPixelBazaVarNotFound, err)
	}

	matches := pixelBazaExpr.FindSubmatch(src)
	if len(matches) < 2 {
		return 0, ErrPixelBazaVarNotFound
	}

	if v, err := strconv.Atoi(string(matches[1])); err == nil {
		return v, nil
	} else {
		return 0, errors.Join(ErrPixelBazaVarNotFound, err)
	}
}

// ParsePixelScriptURL gets the script URL of the pixel challenge script and the URL
// to post a generated payload to from the given HTML code src.
func ParsePixelScriptURL(reader io.Reader) (string, string, error) {
	src, err := io.ReadAll(reader)
	if err != nil {
		return "", "", errors.Join(ErrPixelScriptNotFound, err)
	}

	matches := pixelScriptExpr.FindSubmatch(src)
	if len(matches) < 2 {
		return "", "", errors.Join(ErrPixelScriptNotFound, err)
	}

	scriptUrl := string(matches[1])

	// Create postUrl
	parts := strings.Split(scriptUrl, "/")
	parts[len(parts)-1] = "pixel_" + parts[len(parts)-1]
	postUrl := strings.Join(parts, "/")

	return scriptUrl, postUrl, nil
}

// ParsePixelScriptVar gets the dynamic value from the pixel script
func ParsePixelScriptVar(reader io.Reader) (string, error) {
	src, err := io.ReadAll(reader)
	if err != nil {
		return "", errors.Join(ErrPixelScriptVarNotFound, err)
	}

	index := pixelScriptVarExpr.FindSubmatch(src)
	if len(index) < 2 {
		return "", ErrPixelScriptVarNotFound
	}
	stringIndex, err := strconv.Atoi(string(index[1]))
	if err != nil {
		return "", ErrPixelScriptVarNotFound
	}

	arrayDeclaration := pixelScriptStringArrayExpr.FindSubmatch(src)
	if len(arrayDeclaration) < 2 {
		return "", ErrPixelScriptVarNotFound
	}

	rawStrings := pixelScriptStringsExpr.FindAllSubmatch(arrayDeclaration[1], -1)
	if stringIndex >= len(rawStrings) {
		return "", ErrPixelScriptVarNotFound
	}

	if len(rawStrings[stringIndex]) < 2 {
		return "", ErrPixelScriptVarNotFound
	}

	if v, err := strconv.Unquote(string(rawStrings[stringIndex][1])); err == nil {
		return v, nil
	} else {
		return "", errors.Join(ErrPixelScriptVarNotFound, err)
	}
}