package main

import(
  "errors"
  "path/filepath"
  "bytes"
  "regexp"
  "github.com/hashicorp/vault/api"
  "fmt"
  "encoding/json"
  "gopkg.in/yaml.v2"
)

func recursiveList(client *api.Client, path string) ([]string, error) {
    paths:=[]
    subpaths, err := client.Logical().List(path)
    if err != nil {
		return nil, err
	}
    for _, subpath := range subpaths {
        if string(subpath.Data["type"])=="content"{
            paths = append(paths, path+"/"+string(subpath.Data["data"]))
        } else if string(subpath.Data["type"])=="directory" {
            subpaths, err = getData(client, path+"/"+string(subpath.Data["data"]))
            if err != nil {
                return nil, err
            }
        }
    }
    return paths, nil
}

func editData(data interface{}, input string, output string, passwordLength int) (interface{}, error) {
    lines:= []string{}
    byaml, _:= yaml.Marshal(data)
    var tree yaml.MapSlice
    if err := yaml.Unmarshal(byaml, &tree); err != nil {
        return nil, err
    }
    var buf bytes.Buffer
    if err := render(&buf, tree, ""); err != nil {
        return nil, err
    }
    return ret, nil
}

func render(w io.Writer, tree yaml.MapSlice, prefix string) error {
    for _, branch := range tree {
        key, ok := branch.Key.(string)
        if !ok {
            return errors.Error("unsupported key type: %T", branch.Key)
        }

        prefix := filepath.Join(prefix, key)

        switch x := branch.Value.(type) {
        default:
            return errors.Error("unsupported value type: %T", branch.Value)

        case yaml.MapSlice:
            // recurse
            if err := render(w, x, prefix); err != nil {
                return err
            }
            continue

        // scalar values
        case string:
        case int:
        case float64:
        // ...
        }

        // print scalar
        if _, err := fmt.Fprintf(w, "%s = %v\n", prefix, branch.Value); err != nil {
            return err
        }
    }

    return nil
}

func vaultCopy (client *api.Client, input string, output string, regExp string, passwordLength int){
    paths, err := recursiveList(client, "kv/"+input)
    if err!=nil {
        panic(err)
    }
    b, _ := json.Marshal(data)
	fmt.Println(string(b))
    pat := regexp.MustCompile("^(.*?)"+input+"(.*)$")
    repl:= "${1}"+output+"${2}"
    for _, path := range paths {
        data, err := client.Logical().Read(path)
	    if err != nil {
		    panic(err)
	    }
        editedData := editData(data, input, output, passwordLength)
        outPath := pat.ReplaceAll(path, repl)
        _, err:=client.Logical.Write(outPath, editedData)
	    if err != nil {
		    panic(err)
	    }
    }
}
