package utils

import (
	"bufio"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/dgraph-io/badger"
	log "github.com/sirupsen/logrus"
)

// IsURL checks if input is an URL
func IsURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// AuthorizedURL scans url to verify if there is no unaccepted site inside
func AuthorizedURL(url string) bool {
	file, err := os.Open("lists/banned-urls.list")

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

//AuthorizedText checks if unaccepted words or sentences are in the text
func AuthorizedText(text string) bool {
	file, err := os.Open("lists/banned-text.list")

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
		if strings.Contains(text, line) {
			log.Warnf("An unauthorized text has been provided : %s", text)
			return false
		}
	}
	return true
}

// CheckCategory verifies if a providen category exists
func CheckCategory(category string) bool {
	var knownCategories = []string{"all", "arts", "diy", "films-series", "misc", "music", "nature", "politics-society",
		"science", "sports", "tech"}

	for _, v := range knownCategories {
		if category == v {
			return true
		}
	}
	return false
}

// AntiSpam checks if an ip can post
// It return true if the client can post
func AntiSpam(ip string) bool {
	var canPost bool = true

	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Check if ip is already in badger database
	err = db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			// fmt.Printf("key=%s\n", k)
			if string(k) == ip {
				canPost = false
			}
		}
		return nil
	})

	return canPost
}

// AddIPToAntiSpam adds a ip to waiting list to avoid spamming
func AddIPToAntiSpam(ip string) {
	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Add ip to badger database
	err = db.Update(func(txn *badger.Txn) error {
		// WithTTL parameter is the delay between each post
		e := badger.NewEntry([]byte(ip), []byte("Waiting")).WithTTL(time.Minute * 3)
		err := txn.SetEntry(e)
		return err
	})
}
