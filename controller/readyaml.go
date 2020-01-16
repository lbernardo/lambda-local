package controller

import (
	"encoding/json"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/ghodss/yaml"
)

func ReadYaml(file string) ([]byte, error) {
	content, err := MapServerlessFileIncludes(file)
	if err != nil {
		return nil, err
	}

	cj, err := json.Marshal(content)
	if err != nil {
		return nil, err
	}

	return cj, nil
}

func MapServerlessFileIncludes(file string) (map[string]interface{}, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	contentMap := make(map[string]interface{}, 0)
	err = yaml.Unmarshal(content, &contentMap)
	if err != nil {
		return nil, err
	}
	r, err := regexp.Compile(`\${file\((.+).yml\)}`)
	if err != nil {
		return nil, err
	}
	for i, v := range contentMap {
		s, ok := v.(string)
		if !ok {
			contentMap[i] = v
			continue
		}
		result := r.FindAllString(s, -1)
		if len(result) <= 0 {
			contentMap[i] = s
			continue
		}
		contentMap[i], err = processFiles(result)
		if err != nil {
			return nil, err
		}
	}
	return contentMap, nil
}

func processFiles(files []string) (content map[string]interface{}, err error) {
	content = make(map[string]interface{}, 0)
	for _, file := range files {
		file = strings.Replace(file, "${file(", "", -1)
		file = strings.Replace(file, ")}", "", -1)
		content, err = MapServerlessFileIncludes(file)
		if err != nil {
			return
		}
	}
	return
}
