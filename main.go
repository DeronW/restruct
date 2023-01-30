package main

import (
	"fmt"
	"mss/restruct/pkg/restruct"
	"mss/restruct/pkg/restruct/rule"
)

const s = `
{
	"name": "john",
	"age":47,
	"other":{
		"gender": "male", 
		"colors": ["black", 7]
	}
}
`

const r = `
{
	"action": "json",
	"params": [],
	"return": "json",
	"fields": [{
		"key": "name",
		"type": "string",
		"target": "Name"
	}, {
		"key": "age",
		"type": "number",
		"target": "Age"
	}, {
		"key": "other",
		"type": "json",
		"fields": [{
			"key": "gender",
			"type": "string",
			"extract": {
				"action": "toBool",
				"params": [],
				"return": "bool",
				"target": "g"
			}
		}, {
			"key": "clolors",
			"type": "json",
			"fields": [{
				"key": "0",
				"type" "string",
				"target": "0"
			}, {
				"key": "1",
				"type": "number",
				"target": "1"
			}]
		}] 
	}]
}
`

const r2 = `
{
	"action": "json",
	"params": [],
	"return": "json",
	"fields": [{
		"key": "other",
		"type": "json",
		"fields": [{
			"key": "gender",
			"type": "string",
			"extract": {
				"action": "toBool",
				"params": [],
				"return": "bool",
				"target": "g2"
			}
		}]
	}]
}
`

func main() {
	e, err := rule.Parse(r2)

	if err != nil {
		fmt.Println(err)
	}
	// a, _ := json.Marshal(e)
	// fmt.Println(string(a))
	r, _ := restruct.Restruct(s, e)
	fmt.Println(r)
	fmt.Println(restruct.ToString(r))
}
