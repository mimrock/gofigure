package gofigure

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
)

type JsonFile struct {
	Filename string
}

func NewJsonFile(name string) *JsonFile {
	return &JsonFile {
		Filename: name,
	}
}

func (j JsonFile) ProcessLayer(layer map[string]interface{}, name string, parent *Category) (*Category, error) {
	current := NewCategory(name, parent)

	for k, v := range layer {
		switch v.(type) {
		case string:
			current.Values.Set(k, v.(string))
		case float64:
			current.Values.Set(k, strconv.FormatFloat(v.(float64), 'f', -1, 64))
		case bool:
			b := v.(bool)
			if (b) {
				current.Values.Set(k, "true")
			} else {
				current.Values.Set(k, "false")
			}
		case []interface{}:
			// Lists are pretty complicated to handle. I'll leave this for right now and implement it if required.
			current.Values.Set(k ,"{{list}}")
		case map[string]interface{}:
			c, err := j.ProcessLayer(v.(map[string]interface{}), k, current)
			if err != nil {
				return nil, fmt.Errorf("JsonFile.ProcessLayer: %s", err.Error())
			}
			current.Categories.Set(k, c)
		// According to the json specs, this should only be 'null'.
		// default:
		//		current.Values.Set(k, "{{null}}")
		}
	}
	return current, nil
}

func (j JsonFile) Parse() (*Category, error) {
	data := map[string]interface{}{}
	bytes, err := ioutil.ReadFile(j.Filename)
	if err != nil {
		return nil, fmt.Errorf("JsonFile.Parse: Readfile failed - %s", err.Error())
	}
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, fmt.Errorf("JsonFile.Parse: Unmarshal failed - %s", err.Error())
	}

	root, err := j.ProcessLayer(data, "/", nil)
	if err != nil {
		return nil, fmt.Errorf("JsonFile.Parse: ProcessLayer failed - %s", err.Error())
	}

	return root, nil
}