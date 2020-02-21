package model

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// ReadBaseRules reads the default rules file.
func ReadBaseRules(libDir string) *Rules {
	f, err := os.Open(filepath.Join(libDir, "default_rules.json"))
	if err != nil {
		log.Fatal("couldn't open default rules file")
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal("couldn't read default rules file")
	}
	ret := &Rules{}
	if err = json.Unmarshal(b, &ret); err != nil {
		log.Fatal("couldn't decode default rules file")
	}
	return ret
}
