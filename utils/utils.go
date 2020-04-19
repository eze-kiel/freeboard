package utils

import (
	"bufio"
	"net/url"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

// IsURL checks if input is an URL
func IsURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// AuthorizedURL scans url to verify if there is no unaccepted site inside
func AuthorizedURL(url string) bool {
	file, err := os.Open("../lists/banned-urls.list")

	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var listlines []string

	for scanner.Scan() {
		listlines = append(listlines, scanner.Text())
	}

	file.Close()

	for _, line := range listlines {
		if strings.Contains(url, line) {
			log.Warnf("An unauthorized URL has been provided : %s", url)
			return false
		}
	}
	return true
}
