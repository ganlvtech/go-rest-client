package main

import (
	"io/ioutil"
	"os"

	"github.com/ganlvtech/go-rest-client/protoc-gen-gorestclient/generator"
)

func main() {
	g := generator.New()
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		g.Error(err, "reading input")
	}
	g.Unmarshal(data)
	g.GenerateAllFiles()
	_, err = os.Stdout.Write(g.Marshal())
	if err != nil {
		g.Error(err, "failed to write output proto")
	}
}
