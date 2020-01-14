package resource

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Hash returns the SHA256 hash sum string
// corresponding to the payload addressable
// with the given url, an error, otherwise
func Hash(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", sha256.Sum256(body)), nil
}
