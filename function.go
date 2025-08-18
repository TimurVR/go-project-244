package code

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sort"
	"fmt"
	"gopkg.in/yaml.v3"
)

func Parsing(file string)(map[string]interface{}){
	data, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	var data1 map[string]interface{}
	raz:=filepath.Ext(file)
	switch raz{
	case ".json":
		err = json.Unmarshal([]byte(data), &data1)
		if err != nil {
			log.Fatal(err)
		}
	case ".yml":
		err = yaml.Unmarshal(data, &data1)
		if err != nil {
			log.Fatal(err)
		}
	}
	return data1
}
func  GenDiff(map1,map2 map[string]interface{})(string){
	keys := []string{}
	for key := range map1 {
        keys = append(keys, key)
    }
	for key := range map2 {
		_,exit:=map1[key]
		if !exit{
			keys = append(keys, key)
		} 
    }
	sort.Strings(keys)
	str:="\n{\n"
	for _,key:=range keys{
		value1,exit1:=map1[key]
		value2,exit2:=map2[key]
		if exit1 && exit2 &&value1==value2{
			temp:=fmt.Sprintf("    %s: %v\n", key, value1)
			str=str+temp
		}else if !exit1 && exit2{
			temp:=fmt.Sprintf("  + %s: %v\n", key, value2)
			str=str+temp
		}else if exit1 && !exit2{
			temp:=fmt.Sprintf("  - %s: %v\n", key, value1)
			str=str+temp
		}else{
			temp:=fmt.Sprintf("  - %s: %v\n", key, value1)
			str=str+temp
			temp2:=fmt.Sprintf("  + %s: %v\n", key, value2)
			str=str+temp2
		}
	}
	str+="}\n"
	return str
}
