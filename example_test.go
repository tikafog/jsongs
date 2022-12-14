// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jsongs_test

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/tikafog/jsongs"
)

func ExampleMarshal() {
	type ColorGroup struct {
		ID     int
		Name   string
		Colors []string
	}
	group := ColorGroup{
		ID:     1,
		Name:   "Reds",
		Colors: []string{"Crimson", "Red", "Ruby", "Maroon"},
	}
	b, err := jsongs.Marshal(group)
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(b)
	// Output:
	// {"ID":1,"Name":"Reds","Colors":["Crimson","Red","Ruby","Maroon"]}
}

type MethodColorGroup struct {
	id          int               `json:"ID" json-getter:"ID" json-setter:"SetID"`
	subGroup    *MethodColorGroup `json:"SubGroup,omitempty"`
	_colorArray []string          `json:"ColorArray,omitempty"`
	Name        string
	Colors      []string
}

func (m *MethodColorGroup) ColorArray() []string {
	return m._colorArray
}

func (m *MethodColorGroup) SetColorArray(_colorArray []string) {
	m._colorArray = _colorArray
}

func (m *MethodColorGroup) SubGroup() *MethodColorGroup {
	return m.subGroup
}

func (m *MethodColorGroup) SetSubGroup(subGroup *MethodColorGroup) {
	m.subGroup = subGroup
}

func (m *MethodColorGroup) ID() int {
	return m.id + 1
}

func (m *MethodColorGroup) SetID(id int) {
	m.id = id
}

func ExampleMarshalMethod() {
	group := MethodColorGroup{
		id: 1,
		subGroup: &MethodColorGroup{
			id:     5,
			Name:   "Sder",
			Colors: []string{"Blue", "Green", "Yellow", "Ruby", "Maroon"},
		},
		_colorArray: []string{"Blue", "Green", "Yellow", "Crimson", "Red", "Ruby", "Maroon"},
		Name:        "Reds",
		Colors:      []string{"Crimson", "Red", "Ruby", "Maroon"},
	}
	b, err := jsongs.Marshal(group)
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(b)
	// Output:
	// {"ID":2,"SubGroup":{"ID":6,"Name":"Sder","Colors":["Blue","Green","Yellow","Ruby","Maroon"]},"ColorArray":["Blue","Green","Yellow","Crimson","Red","Ruby","Maroon"],"Name":"Reds","Colors":["Crimson","Red","Ruby","Maroon"]}
}

func ExampleUnmarshal() {
	var jsonBlob = []byte(`[
	{"Name": "Platypus", "Order": "Monotremata"},
	{"Name": "Quoll",    "Order": "Dasyuromorphia"}
]`)
	type Animal struct {
		Name  string
		Order string
	}
	var animals []Animal
	err := jsongs.Unmarshal(jsonBlob, &animals)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("%+v", animals)
	// Output:
	// [{Name:Platypus Order:Monotremata} {Name:Quoll Order:Dasyuromorphia}]
}

type MethodSubAnimal struct {
	name  string `json:"name" json-setter:"MyName"`
	order string `json:"order" json-getter:"MyOrder"`
}

func (m *MethodSubAnimal) Name() string {
	return m.name
}

func (m *MethodSubAnimal) MyOrder() string {
	return m.order
}

func (m *MethodSubAnimal) SetOrder(order string) {
	m.order = order
}

func (m *MethodSubAnimal) MyName(name string) {
	m.name = name
}

type MethodAnimal struct {
	name      string          `json:"name" json-getter:"InName"`
	Name      string          `json:"Name"`
	order     string          `json:"order" json-getter:"Order"`
	subAnimal MethodSubAnimal `json:"sub_animal"`
	//Order string
}

func (receiver *MethodAnimal) InName() string {
	return receiver.name
}

func (receiver *MethodAnimal) SubAnimal() MethodSubAnimal {
	return receiver.subAnimal
}

func (receiver *MethodAnimal) SetSubAnimal(subAnimal MethodSubAnimal) {
	receiver.subAnimal = subAnimal
}

func (receiver *MethodAnimal) Order() string {
	return receiver.order
}

func (receiver *MethodAnimal) SetOrder(order string) {
	receiver.order = order
}

func (receiver *MethodAnimal) SetName(name string) {
	receiver.name = name
}

func (receiver *MethodAnimal) MyName(name string) {
	receiver.name = name
}

func ExampleUnmarshalMethod() {
	var jsonBlob = []byte(`[
	{"name": "myname", "sub_animal": {"Name": "Quoll",    "Order": "Dasyuromorphia"}, "Name": "Platypus", "Order": "Monotremata"},
	{"name": "ismyname", "sub_animal":{"Name": "Platypus", "Order": "Monotremata"}, "Name": "Quoll",    "Order": "Dasyuromorphia"}
]`)

	var animals []MethodAnimal
	err := jsongs.Unmarshal(jsonBlob, &animals)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("%+v\n", animals)
	fmt.Printf("%+v\n", animals[0].subAnimal)
	fmt.Printf("%+v", animals[1].subAnimal)
	// Output:
	// [{name:myname Name:Platypus order:Monotremata subAnimal:{name:Quoll order:Dasyuromorphia}} {name:ismyname Name:Quoll order:Dasyuromorphia subAnimal:{name:Platypus order:Monotremata}}]
	// {name:Quoll order:Dasyuromorphia}
	// {name:Platypus order:Monotremata}
}

// This example uses a Decoder to decode a stream of distinct JSON values.
func ExampleDecoder() {
	const jsonStream = `
	{"Name": "Ed", "Text": "Knock knock."}
	{"Name": "Sam", "Text": "Who's there?"}
	{"Name": "Ed", "Text": "Go fmt."}
	{"Name": "Sam", "Text": "Go fmt who?"}
	{"Name": "Ed", "Text": "Go fmt yourself!"}
`
	type Message struct {
		Name, Text string
	}
	dec := jsongs.NewDecoder(strings.NewReader(jsonStream))
	for {
		var m Message
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s: %s\n", m.Name, m.Text)
	}
	// Output:
	// Ed: Knock knock.
	// Sam: Who's there?
	// Ed: Go fmt.
	// Sam: Go fmt who?
	// Ed: Go fmt yourself!
}

// This example uses a Decoder to decode a stream of distinct JSON values.
func ExampleDecoder_Token() {
	const jsonStream = `
	{"Message": "Hello", "Array": [1, 2, 3], "Null": null, "Number": 1.234}
`
	dec := jsongs.NewDecoder(strings.NewReader(jsonStream))
	for {
		t, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%T: %v", t, t)
		if dec.More() {
			fmt.Printf(" (more)")
		}
		fmt.Printf("\n")
	}
	// Output:
	// jsongs.Delim: { (more)
	// string: Message (more)
	// string: Hello (more)
	// string: Array (more)
	// jsongs.Delim: [ (more)
	// float64: 1 (more)
	// float64: 2 (more)
	// float64: 3
	// jsongs.Delim: ] (more)
	// string: Null (more)
	// <nil>: <nil> (more)
	// string: Number (more)
	// float64: 1.234
	// jsongs.Delim: }
}

// This example uses a Decoder to decode a streaming array of JSON objects.
func ExampleDecoder_Decode_stream() {
	const jsonStream = `
	[
		{"Name": "Ed", "Text": "Knock knock."},
		{"Name": "Sam", "Text": "Who's there?"},
		{"Name": "Ed", "Text": "Go fmt."},
		{"Name": "Sam", "Text": "Go fmt who?"},
		{"Name": "Ed", "Text": "Go fmt yourself!"}
	]
`
	type Message struct {
		Name, Text string
	}
	dec := jsongs.NewDecoder(strings.NewReader(jsonStream))

	// read open bracket
	t, err := dec.Token()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%T: %v\n", t, t)

	// while the array contains values
	for dec.More() {
		var m Message
		// decode an array value (Message)
		err := dec.Decode(&m)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%v: %v\n", m.Name, m.Text)
	}

	// read closing bracket
	t, err = dec.Token()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%T: %v\n", t, t)

	// Output:
	// jsongs.Delim: [
	// Ed: Knock knock.
	// Sam: Who's there?
	// Ed: Go fmt.
	// Sam: Go fmt who?
	// Ed: Go fmt yourself!
	// jsongs.Delim: ]
}

// This example uses RawMessage to delay parsing part of a JSON message.
func ExampleRawMessage_unmarshal() {
	type Color struct {
		Space string
		Point jsongs.RawMessage // delay parsing until we know the color space
	}
	type RGB struct {
		R uint8
		G uint8
		B uint8
	}
	type YCbCr struct {
		Y  uint8
		Cb int8
		Cr int8
	}

	var j = []byte(`[
	{"Space": "YCbCr", "Point": {"Y": 255, "Cb": 0, "Cr": -10}},
	{"Space": "RGB",   "Point": {"R": 98, "G": 218, "B": 255}}
]`)
	var colors []Color
	err := jsongs.Unmarshal(j, &colors)
	if err != nil {
		log.Fatalln("error:", err)
	}

	for _, c := range colors {
		var dst any
		switch c.Space {
		case "RGB":
			dst = new(RGB)
		case "YCbCr":
			dst = new(YCbCr)
		}
		err := jsongs.Unmarshal(c.Point, dst)
		if err != nil {
			log.Fatalln("error:", err)
		}
		fmt.Println(c.Space, dst)
	}
	// Output:
	// YCbCr &{255 0 -10}
	// RGB &{98 218 255}
}

// This example uses RawMessage to use a precomputed JSON during marshal.
func ExampleRawMessage_marshal() {
	h := jsongs.RawMessage(`{"precomputed": true}`)

	c := struct {
		Header *jsongs.RawMessage `json:"header"`
		Body   string             `json:"body"`
	}{Header: &h, Body: "Hello Gophers!"}

	b, err := jsongs.MarshalIndent(&c, "", "\t")
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(b)

	// Output:
	// {
	// 	"header": {
	// 		"precomputed": true
	// 	},
	// 	"body": "Hello Gophers!"
	// }
}

func ExampleIndent() {
	type Road struct {
		Name   string
		Number int
	}
	roads := []Road{
		{"Diamond Fork", 29},
		{"Sheep Creek", 51},
	}

	b, err := jsongs.Marshal(roads)
	if err != nil {
		log.Fatal(err)
	}

	var out bytes.Buffer
	jsongs.Indent(&out, b, "=", "\t")
	out.WriteTo(os.Stdout)
	// Output:
	// [
	// =	{
	// =		"Name": "Diamond Fork",
	// =		"Number": 29
	// =	},
	// =	{
	// =		"Name": "Sheep Creek",
	// =		"Number": 51
	// =	}
	// =]
}

func ExampleMarshalIndent() {
	data := map[string]int{
		"a": 1,
		"b": 2,
	}

	b, err := jsongs.MarshalIndent(data, "<prefix>", "<indent>")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))
	// Output:
	// {
	// <prefix><indent>"a": 1,
	// <prefix><indent>"b": 2
	// <prefix>}
}

func ExampleValid() {
	goodJSON := `{"example": 1}`
	badJSON := `{"example":2:]}}`

	fmt.Println(jsongs.Valid([]byte(goodJSON)), jsongs.Valid([]byte(badJSON)))
	// Output:
	// true false
}

func ExampleHTMLEscape() {
	var out bytes.Buffer
	jsongs.HTMLEscape(&out, []byte(`{"Name":"<b>HTML content</b>"}`))
	out.WriteTo(os.Stdout)
	// Output:
	//{"Name":"\u003cb\u003eHTML content\u003c/b\u003e"}
}

type Example struct {
	name string `json:"name" json-getter:"MyName" json-setter:"SetMyName"`
}

func (receiver Example) MyName() string {
	return receiver.name
}

func (receiver *Example) SetMyName(name string) {
	receiver.name = name
}

func ExampleMainFunction() {
	v, err := jsongs.Marshal(&Example{
		name: "my name is jsongs",
	})
	if err != nil {
		panic(err)
	}
	//dosomething
	fmt.Println(string(v))
	// Output:
	//{"name":"my name is jsongs"}
}
