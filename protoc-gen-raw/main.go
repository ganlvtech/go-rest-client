package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

func main() {
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err, "reading input")
	}

	response := &plugin.CodeGeneratorResponse{}
	response.File = append(response.File, &plugin.CodeGeneratorResponse_File{
		Name:    proto.String("raw.protobuf"),
		Content: proto.String(string(data)),
	})

	data, err = proto.Marshal(response)
	if err != nil {
		log.Fatal(err, "failed to marshal output proto")
	}

	_, err = os.Stdout.Write(data)
	if err != nil {
		log.Fatal(err, "failed to write output proto")
	}
}
