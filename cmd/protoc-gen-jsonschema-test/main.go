package main

import (
	"io"
	"log"
	"os"
	"strings"

	plugingo "github.com/golang/protobuf/protoc-gen-go/plugin"
	"google.golang.org/protobuf/proto"
)

func main() {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal("unable to read input: ", err)
	}

	req := &plugingo.CodeGeneratorRequest{}
	if err = proto.Unmarshal(data, req); err != nil {
		log.Fatal("unable to unmarshal request: ", err)
	}
	toGenerate := req.GetFileToGenerate()
	switch len(toGenerate) {
	case 0:
		return
	case 1:
		break
	default:
		log.Fatal("this tool only support one file per call")
	}
	sourceFile := toGenerate[0]
	binFile := strings.TrimSuffix(sourceFile, "proto") + "pb.bin"
	file, err := os.OpenFile(binFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	defer func() { _ = file.Close() }()
	if err != nil {
		log.Fatalf("create bin file failed: %+v", err)
	}
	binContent, err := proto.Marshal(req)
	if err != nil {
		log.Fatalf("error marshalling request: %+v", err)
	}
	_, err = file.Write(binContent)
	if err != nil {
		log.Fatalf("error saving request: %+v", err)
	}
}
