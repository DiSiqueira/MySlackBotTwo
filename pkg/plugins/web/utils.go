package web

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"net/url"
)

const (
	minDomainLength = 3
)

func GetBody(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}

func GetJSON(url string, v interface{}) error {
	body, err := GetBody(url)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, v)
}

func canBeURLWithoutProtocol(text string) bool {
	return len(text) > minDomainLength &&
		!strings.HasPrefix(text, "http") &&
		strings.Contains(text, ".")
}

func ExtractURL(text string) string {
	extractedURL := ""
	for _, value := range strings.Split(text, " ") {
		value = cleanURL(value, []string{" ", ">", "<"})
		if canBeURLWithoutProtocol(value) {
			value = "http://" + value
		}

		parsedURL, err := url.Parse(value)
		if err != nil {
			continue
		}
		if strings.HasPrefix(parsedURL.Scheme, "http") {
			extractedURL = parsedURL.String()
			break
		}
	}
	return extractedURL
}


func cleanURL(url string, cutset []string) string {
	final := url
	for _, cut := range cutset {
		final = strings.Trim(final, cut)
	}

	return final
}