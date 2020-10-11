// +build generate

package main

import (
	"encoding/base64"
	"io/ioutil"
)

func main() {
	b, err := ioutil.ReadFile("./template.gotemplate")
	if err != nil {
		panic(err)
	}
	encodedtmpl := base64.StdEncoding.EncodeToString(b)

	b = []byte(`// Code generated by generate.go DO NOT EDIT.
package style

import (
	"encoding/base64"
	"text/template"
)

const (
        raw = ` + "`" + encodedtmpl + "`" + `
)

var (
	Template = template.Must(template.New("").Parse(mustString(base64.StdEncoding.DecodeString(raw))))
)

func mustString(b []byte, err error) string {
	if err != nil {
		panic(err)
	}
	return string(b)
}
`)

	err = ioutil.WriteFile("style.go", b, 0644)
	if err != nil {
		panic(err)
	}
}