package plugin

import (
	"fmt"
	"strings"

	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
)

// define the plugin struct
type plugin struct {
	*generator.Generator
	generator.PluginImports
	protoPkg generator.Single
}

// New will return a new plugin instance
func New() generator.Plugin {
	return &plugin{}
}

// Name of the plugin
func (p *plugin) Name() string {
	return "gotags"
}

// Init will initialze the plugin with the given generator
func (p *plugin) Init(g *generator.Generator) {
	p.Generator = g
}

func (p *plugin) Generate(file *generator.FileDescriptor) {

	// we use a custom struct to handle package imports
	p.PluginImports = NewPluginImports(p.Generator)

	// extract all comments from the file into a map with the path
	// to the comment as key
	comments := extractComments(file)

	// initialize a map for the message path, message paths are required
	// to find comments for the message and fields
	pathIndex := make(map[string]string)

	// generate the path index
	for mIndex, message := range file.Messages() {
		// 4 at the beginning is for type message
		pathIndex[message.GetName()] = fmt.Sprintf("4,%d", mIndex)
	}

	// generate custom structs
	generateStructs(p, file, comments, pathIndex)

	// NOTE: conversion functions are killed for now. it's just to complicated
	// p.P("\n // --- STRUCT CONVERSION --- \n")

	// // generate functions to convert structs
	// generateConvertFnCustom(p, file, comments, pathIndex)
	// generateConvertFnProto(p, file, comments, pathIndex)

}

// Printf will allow us to use printf syntax for printing
func (p *plugin) Printf(layout string, args ...interface{}) {
	p.P(fmt.Sprintf(layout, args...))
}

// PrintReturnErr will print an error handler
func (p *plugin) PrintReturnErr(text ...string) {

	p.P("if err != nil {")
	p.In()

	if len(text) > 0 {
		p.P(strings.Join(text, ""))
	} else {
		p.P("return nil, err")
	}

	p.Out()
	p.P("}")

}
