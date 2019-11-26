package controller

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
)

func ReadYaml(file string) ([]byte, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	cj, err := yaml.YAMLToJSON(content)
	if err != nil {
		return nil, err
	}

	return cj, nil
}
