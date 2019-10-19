package chglog

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

func Parse(file string) (entries ChangeLogEntries, err error) {
	var (
		body []byte
	)
	body, err = ioutil.ReadFile(file)
	switch {
	case os.IsNotExist(err):
		return make(ChangeLogEntries, 0), nil
	case err != nil:
		return nil, err
	}

	err = yaml.Unmarshal(body, &entries)
	return entries, err
}

func (cle *ChangeLogEntries) Save(file string) (err error) {
	data, _ := yaml.Marshal(cle)
	return ioutil.WriteFile(file, data, 0644)
}
