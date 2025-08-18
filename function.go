package code

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"gopkg.in/yaml.v3"
)

func Parsing(file string){
	data, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	var data1 map[string]interface{}
	raz:=filepath.Ext(file)
	if raz==".json"{
		err = json.Unmarshal([]byte(data), &data1)
		if err != nil {
			log.Fatal(err)
		}
	}else if raz==".yml"{
		err = yaml.Unmarshal(data, &data1)
		if err != nil {
			log.Fatal(err)
		}

	}
	fmt.Println(data1)
}
