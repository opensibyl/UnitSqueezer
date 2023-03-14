package main

import (
	"encoding/json"
	"flag"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/opensibyl/squ"
	"github.com/opensibyl/squ/object"
)

const ConfigFile = "squ.config.json"

func main() {
	config := object.DefaultConfig()

	// cmd parse
	src := flag.String("src", ".", "repo path")
	before := flag.String("before", "HEAD~1", "before rev")
	after := flag.String("after", "HEAD", "after rev")
	diffOutput := flag.String("jsonOutput", "", "diff output")
	graphOutput := flag.String("graphOutput", "", "svg output")
	dry := flag.Bool("dry", false, "dry")
	indexerType := flag.String("indexer", object.IndexerGolang, "indexer type")
	runnerType := flag.String("runner", object.RunnerGolang, "runner type")
	sibylUrl := flag.String("sibylUrl", config.SibylUrl, "url of sibyl server")
	debugMode := flag.Bool("debug", config.DebugMode, "debug mode switch")
	flag.Parse()

	// load data from config file
	configFile := filepath.Join(*src, ConfigFile)
	if _, err := os.Stat(configFile); err == nil {
		data, err := os.ReadFile(configFile)
		squ.PanicIfErr(err)
		err = json.Unmarshal(data, &config)
		squ.PanicIfErr(err)
	}

	config.SrcDir = *src
	config.Before = *before
	config.After = *after
	config.JsonOutput = *diffOutput
	config.GraphOutput = *graphOutput
	config.Dry = *dry
	config.IndexerType = *indexerType
	config.RunnerType = *runnerType
	config.SibylUrl = *sibylUrl
	config.DebugMode = *debugMode

	bytes, err := json.MarshalIndent(config, "", "    ")
	squ.PanicIfErr(err)
	err = os.WriteFile(configFile, bytes, fs.ModePerm)
	squ.PanicIfErr(err)

	squ.MainFlow(config)
}
