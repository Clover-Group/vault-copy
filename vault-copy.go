package main

import (
    "strings"
    pureJson "encoding/json"
	"fmt"
	"github.com/hashicorp/vault/api"
	"gopkg.in/yaml.v2"
    "github.com/tidwall/sjson"
	"regexp"
)

func recursiveList(client *api.Client, path string) ([]string, error) {
	paths := []string{}
	subpaths, err := client.Logical().List("kv/metadata/" + path)
	if err != nil {
		return nil, err
	}
	rawKeys := subpaths.Data["keys"]
	for _, key := range rawKeys.([]interface{}) {
		subpath := key.(string)
		if subpath[len(subpath)-1] != '/' {
			paths = append(paths, path+subpath)
		} else {
			tmp, err := recursiveList(client, path+subpath)
			if err != nil {
				return nil, err
			}
			paths = append(paths, tmp...)
		}
	}
	return paths, nil
}

func editData(data interface{}, input string, output string, passwordLength int) (string, error) {
	byaml, _ := yaml.Marshal(data)
	var tree yaml.MapSlice
	if err := yaml.Unmarshal(byaml, &tree); err != nil {
		return "", err
	}
	lines := map[string]interface{}{}
	if err := plain(lines, tree, ""); err != nil {
		return "", err
	}
	pat := regexp.MustCompile("^(.*?)" + input + "(.*)$")
	repl := "${1}" + output + "${2}"
    json:=""
    var err error
    for k, v := range lines {
        if strings.Contains(k, "password") {
            lines[k]=randomString(passwordLength)
        }
        if v1, ok:=v.(string); ok{
            if strings.Contains(v1, input) {
		        out := string(pat.ReplaceAll([]byte(v1), []byte(repl)))
                lines[k]=out
            }
        }
        json, err = sjson.Set(json, k, lines[k])
        if err!=nil {
            return "", err
        }
    }
	return json, nil
}

func plain(lines map[string]interface{}, tree yaml.MapSlice, prefix string) error {
	for _, branch := range tree {
		key, ok := branch.Key.(string)
		if !ok {
			return fmt.Errorf("unsupported key type: %T", branch.Key)
		}
		newPrefix := ""
		if prefix != "" {
			newPrefix = prefix + "." + key
		} else {
			newPrefix = key
		}

		switch x := branch.Value.(type) {
		default:
			return fmt.Errorf("unsupported value type: %T", branch.Value)
		case yaml.MapSlice:
			// recurse
			if err := plain(lines, x, newPrefix); err != nil {
				return err
			}
			continue
		case []interface{}:
		case string:
		case int:
		case bool:
		case float64:
		case nil:
		}
		lines[newPrefix] = branch.Value
	}

	return nil
}

func vaultCopy(client *api.Client, input string, output string, regExp string, passwordLength int) {
	paths, err := recursiveList(client, input+"/")
	if err != nil {
		panic(err)
	}
	pat := regexp.MustCompile("^(.*?)" + input + "(.*)$")
	repl := "${1}" + output + "${2}"
	for _, path := range paths {
		data, err := client.Logical().Read("kv/data/" + path)
		if err != nil {
			panic(err)
		}
		editedData, err := editData(data.Data["data"], input, output, passwordLength)
		if err != nil {
			panic(err)
		}
		outPath := string(pat.ReplaceAll([]byte(path), []byte(repl)))
        var b1 map[string]interface{}
        pureJson.Unmarshal([]byte(editedData), &b1)
        b:=map[string]interface{}{"data": b1}
		_, err = client.Logical().Write("kv/data/"+outPath, b)
		if err != nil {
			panic(err)
		}
	}
}
