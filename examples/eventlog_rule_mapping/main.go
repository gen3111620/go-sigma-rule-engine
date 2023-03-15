package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/markuskont/datamodels"
	"github.com/markuskont/go-sigma-rule-engine"
)

var (
	flagRuleSetPath = flag.String("path-ruleset", "./windows/", "Root folders for Sigma rules. Semicolon delimits paths.")
)

func main() {
	flag.Parse()
	if *flagRuleSetPath == "" {
		log.Fatal("ruleset path not configured")
	}
	ruleset, err := sigma.NewRuleset(sigma.Config{
		Directory:       strings.Split(*flagRuleSetPath, ";"),
		NoCollapseWS:    false,
		FailOnRuleParse: false,
		FailOnYamlParse: false,
	})
	if err != nil {
		log.Fatal(err)
	}
	data, err := ioutil.ReadFile("data.json")

	if err != nil {
		log.Fatal(err)
	}

	var events []map[string]interface{}

	if err := json.Unmarshal([]byte(data), &events); err != nil {
		panic(err)
	}

	output := os.Stdout
	for _, event := range events {
		jsonStr, err := json.Marshal(event)
		if err != nil {
			log.Println(err)
		}

		var obj datamodels.Map
		if err := json.Unmarshal(jsonStr, &obj); err != nil {
			log.Println(err)
		}

		if results, ok := ruleset.EvalAll(obj); ok && len(results) > 0 {
			obj["sigma_results"] = results
			encoded, err := json.Marshal(obj)
			if err != nil {
				log.Println(err)
			}
			output.Write(encoded)
		}

	}

}
