package services

import (
	"crypto/rand"
	"strings"
	"testing"

	uuid "github.com/satori/go.uuid"
)

func Assert(t *testing.T, fail bool) func(bool, ...interface{}) {
	if fail {
		return func(expression bool, output ...interface{}) {
			if !expression {
				t.Fatal(output)
			}
		}
	} else {
		return func(expression bool, output ...interface{}) {
			if !expression {
				t.Log(output)
			}
		}
	}
}

func GenerateId() string {
	b := make([]byte, 16)
	rand.Read(b)
	newId, _ := uuid.NewV4()
	return newId.String()
}

// ValidURL if count '/' in url is more than expected returns false
func ValidURL(url string, count int) bool {
	if strings.Count(url, "/") > count {
		return false
	}
	return true
}
