package code
import (
	"log"
	"os"
)

func Parsing(){
	data, err := os.ReadFile("testdata/hello")
	if err != nil {
		log.Fatal(err)
	}
	os.Stdout.Write(data)
}
