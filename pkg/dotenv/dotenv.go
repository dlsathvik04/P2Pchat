package dotenv

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func LoadDotEnv(filename string, soft bool) {
	fmt.Println("### configureing env ###########")
	data, err := os.ReadFile(filename)
	if err != nil {
		if soft {
			return
		}
		log.Fatal(err)
	}
	datastr := string(data)
	lines := strings.Split(datastr, "\n")
	for _, line := range lines {
		if !(strings.HasPrefix(line, "#") || strings.Trim(line, " ") == "") {
			currpair := strings.SplitN(line, "=", 2)
			err := os.Setenv(currpair[0], currpair[1])
			fmt.Println("Set ", currpair[0], " -> ", currpair[1])
			if err != nil {
				log.Fatal("Cannot set env")
			}
		}
	}
	fmt.Println("configured env from : ", filename)

	fmt.Println("#################################")
}
