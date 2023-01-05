package main

import (
	"os"
	"strings"
	"testing"

	"golang.org/x/term"
)

func TestMatchesAnyDomain(t *testing.T) {
	doesMatch := matchesAnyDomain("cloud.reallyawesomedomain.com", []string{"cloud", "foo"})
	doesNotMatch := matchesAnyDomain("bar.example.com", []string{"cloud", "foo"})
	if !doesMatch {
		t.Error("Expected to match got:", doesMatch)
	}
	if doesNotMatch {
		t.Error("Expected not to match got:", doesMatch)
	}
}

func TestParseDomainConf(t *testing.T) {
	parsedMap := parseDomainConf("foo.com:cloud,bar.net:www")
	if len(parsedMap) < 2 {
		t.Error("Expected parsedMap to have 2 entries got:", parsedMap)
	}
	if parsedMap["foo.com"][0] != "cloud" || parsedMap["bar.net"][0] != "www" {
		t.Errorf("Expected a correct parse go: %T", parsedMap)
	}

	parsedMapWithMultipleDomains := parseDomainConf("foo.com:cloud,foo.com:www")
	if len(parsedMapWithMultipleDomains["foo.com"]) < 2 {
		t.Error("Expected foo.com record to have 2 entries got:", parsedMapWithMultipleDomains["foo.com"])
	}
	if parsedMapWithMultipleDomains["foo.com"][0] != "cloud" || parsedMapWithMultipleDomains["foo.com"][1] != "www" {
		t.Error("Expected a correct parse got:", parsedMapWithMultipleDomains)
	}
}

func TestAskAHumanFrom(t *testing.T) {
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		t.Skip("Not running in a terminal. Run `go test -run TestAskAHumanFrom` in a terminal to test")
	}
	readFromYes := strings.NewReader("y\n")
	yes := askAHumanFrom(readFromYes, "test")
	if !yes {
		t.Error("Expected a yes response")
	}
	readFromNo := strings.NewReader("n\n")
	no := askAHumanFrom(readFromNo, "test")
	if no {
		t.Error("Expected a no response")
	}
}
