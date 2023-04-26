//go:build ignore
// +build ignore

// Copyright 2017 Google LLC. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// This is a helper utility to embed the Redis Lua scripts into a Go source
// file.
package main

import (
	"bytes"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"text/template"
)

var packageTemplate = template.Must(template.New("").Parse(`
// Code generated by quota/redis/redistb/gen.go. DO NOT EDIT.
// source: {{ .Filename }}

package redistb

import (
	"github.com/go-redis/redis"
)

// contents of the '{{ .Prefix }}' Redis Lua script
const {{ .Prefix }}ScriptContents = {{ .Content }}

// Redis Script type for the '{{ .Prefix }}' Redis lua script
var {{ .Prefix }}Script = redis.NewScript({{ .Prefix }}ScriptContents)
`))

type templateData struct {
	Prefix   string
	Filename string
	Content  string
}

func main() {
	if len(os.Args) != 4 {
		log.Fatalf("usage: %s prefix file.lua output.go", os.Args[0])
	}

	data, err := ioutil.ReadFile(os.Args[2])
	if err != nil {
		log.Fatalf("error opening input file: %v", err)
	}

	vars := templateData{
		Prefix:   os.Args[1],
		Filename: os.Args[2],
		Content:  strconv.Quote(string(data)),
	}

	var buf bytes.Buffer
	if err := packageTemplate.Execute(&buf, vars); err != nil {
		log.Fatalf("error rendering template: %v", err)
	}

	data, err = format.Source(buf.Bytes())
	if err != nil {
		log.Fatalf("error formatting source: %v", err)
	}

	out, err := os.Create(os.Args[3])
	if err != nil {
		log.Fatalf("error opening output file: %v", err)
	}
	defer out.Close()

	if _, err := out.Write(data); err != nil {
		log.Fatalf("error writing output file: %v", err)
	}
}
