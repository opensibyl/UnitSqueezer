package squ

import (
	"fmt"
	"path"
	"path/filepath"
	"testing"

	"github.com/opensibyl/squ/object"
)

func TestGolang(t *testing.T) {
	conf := object.DefaultConfig()
	conf.Dry = true
	conf.GraphOutput = "abc.png"
	conf.DebugMode = true
	MainFlow(conf)
}

func TestMaven(t *testing.T) {
	t.Skip()
	conf := object.DefaultConfig()
	conf.Dry = true
	conf.SrcDir = "../jacoco"
	conf.Before = "HEAD~5"
	conf.IndexerType = object.IndexerJavaJUnit
	conf.RunnerType = object.RunnerMaven
	conf.JsonOutput = "./output.json"
	MainFlow(conf)
}

func TestA(t *testing.T) {
	files, _ := filepath.Glob(path.Join("./stagesepx", "*"))
	for _, each := range files {
		fmt.Println(each)
	}
}
