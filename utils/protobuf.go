package utils

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/golang/protobuf/protoc-gen-go/plugin"

	"errors"
	"io/ioutil"
	"path/filepath"
)

// CreateGenRequest creates a codegen request from a `FileDescriptorSet`
func CreateGenRequest(fds *descriptor.FileDescriptorSet, filesToGen ...string) *plugin_go.CodeGeneratorRequest {
	req := new(plugin_go.CodeGeneratorRequest)
	req.ProtoFile = fds.GetFile()

	for _, f := range req.GetProtoFile() {
		if InStringSlice(filesToGen, filepath.Base(f.GetName())) {
			req.FileToGenerate = append(req.FileToGenerate, f.GetName())
		}
	}

	return req
}

// LoadDescriptorSet loads a `FileDescriptorSet` from a file on disk. Such a file can be generated using the
// `--descriptor_set_out` flag with `protoc`.
//
// Example:
//     protoc --descriptor_set_out=fileset.pb --include_imports --include_source_info ./booking.proto ./todo.proto
func LoadDescriptorSet(pathSegments ...string) (*descriptor.FileDescriptorSet, error) {
	f, err := ioutil.ReadFile(filepath.Join(pathSegments...))
	if err != nil {
		return nil, err
	}

	set := new(descriptor.FileDescriptorSet)
	if err = proto.Unmarshal(f, set); err != nil {
		return nil, err
	}

	return set, nil
}

// FindDescriptor finds the named descriptor in the given set. Only base names are searched. The first match is
// returned, on `nil` if not found
func FindDescriptor(set *descriptor.FileDescriptorSet, name string) *descriptor.FileDescriptorProto {
	for _, pf := range set.GetFile() {
		if filepath.Base(pf.GetName()) == name {
			return pf
		}
	}

	return nil
}

// LoadDescriptor loads file descriptor protos from a file on disk, and returns the named proto descriptor. This is
// useful mostly for testing purposes.
func LoadDescriptor(name string, pathSegments ...string) (*descriptor.FileDescriptorProto, error) {
	set, err := LoadDescriptorSet(pathSegments...)
	if err != nil {
		return nil, err
	}

	if pf := FindDescriptor(set, name); pf != nil {
		return pf, nil
	}

	return nil, errors.New("FileDescriptor not found")
}
