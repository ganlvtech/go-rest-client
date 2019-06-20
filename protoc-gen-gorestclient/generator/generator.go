package generator

import (
	"log"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

type Generator struct {
	Request  *plugin.CodeGeneratorRequest
	Response *plugin.CodeGeneratorResponse

	fileByName map[string]*descriptor.FileDescriptorProto

	FileGenerators   []*FileGenerator
	FileGeneratorMap map[string]*FileGenerator
}

func New() *Generator {
	g := new(Generator)
	g.Request = new(plugin.CodeGeneratorRequest)
	g.Response = new(plugin.CodeGeneratorResponse)
	g.FileGeneratorMap = make(map[string]*FileGenerator)
	return g
}

// Error reports a problem, including an error, and exits the program.
func (g *Generator) Error(err error, msgs ...string) {
	s := strings.Join(msgs, " ") + ":" + err.Error()
	log.Fatal("protoc-gen-go: error:", s)
	// unreachable
}

// Fail reports a problem and exits the program.
func (g *Generator) Fail(msgs ...string) {
	s := strings.Join(msgs, " ")
	log.Fatal("protoc-gen-go: error:", s)
	// unreachable
}

func (g *Generator) Unmarshal(data []byte) {
	if err := proto.Unmarshal(data, g.Request); err != nil {
		g.Error(err, "parsing input proto")
	}
}

func (g *Generator) Marshal() []byte {
	data, err := proto.Marshal(g.Response)
	if err != nil {
		g.Error(err, "failed to marshal output proto")
	}
	return data
}

func (g *Generator) mapFileByName() {
	g.fileByName = make(map[string]*descriptor.FileDescriptorProto)
	for _, file := range g.Request.ProtoFile {
		g.fileByName[file.GetName()] = file
	}
}

func (g *Generator) GenerateAllFiles() {
	if len(g.Request.FileToGenerate) == 0 {
		g.Fail("no files to generate")
	}
	g.mapFileByName()
	for _, fileName := range g.Request.FileToGenerate {
		g.GenerateFile(g.fileByName[fileName])
	}
}

func (g *Generator) GenerateFile(file *descriptor.FileDescriptorProto) {
	fileGenerator := NewFileGenerator(file)
	fileGenerator.WrapTypes()
	fileGenerator.Generate()
	g.Response.File = append(g.Response.File, &plugin.CodeGeneratorResponse_File{
		Name:    proto.String(fileGenerator.GoFileName()),
		Content: proto.String(fileGenerator.String()),
	})
}
