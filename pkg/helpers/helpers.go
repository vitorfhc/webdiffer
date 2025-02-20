package helpers

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"
)

func ReadOrCreateJSON(filepath string) ([]byte, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			err = CreateEmptyJSON(filepath)
			if err != nil {
				return nil, err
			}
			return ReadOrCreateJSON(filepath)
		}
		return nil, err
	}

	return data, nil
}

func CreateEmptyJSON(filepath string) error {
	emptyData, err := json.Marshal([]interface{}{})
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath, emptyData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func NormalizeURL(input string) (string, error) {
	parsedURL, err := url.Parse(input)
	if err != nil || parsedURL.Host == "" {
		return "", fmt.Errorf("invalid URL: %s", input)
	}

	if parsedURL.Scheme == "" {
		parsedURL.Scheme = "https"
	}

	parsedURL.Host = strings.ToLower(parsedURL.Host)

	if (parsedURL.Scheme == "http" && strings.HasSuffix(parsedURL.Host, ":80")) ||
		(parsedURL.Scheme == "https" && strings.HasSuffix(parsedURL.Host, ":443")) {
		parsedURL.Host = strings.Split(parsedURL.Host, ":")[0]
	}

	normalizedPath := strings.TrimRight(parsedURL.Path, "/")

	return fmt.Sprintf("%s://%s%s", parsedURL.Scheme, parsedURL.Host, normalizedPath), nil
}
