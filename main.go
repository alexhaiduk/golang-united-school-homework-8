package main

import (
	"io"
	//"io/ioutil"
	"os"
	//"bytes"
	"flag"
	//"errors"
)

type Arguments map[string]string

func parseArgs() Arguments {
	var args Arguments
	args["fileName"] = *flag.String("fileName", "users.json", "name of file where you store users, default: \"users.json\"")
	args["id"] = *flag.String("id", "", "user id, default: \"\"")
	args["item"] = *flag.String("item", "", "valid json object with the id, email and age fields") //maybe default: {\"id\": \"0\", \"email\": \"email@mail.com\", \"age\": 1}
	args["operation"] = *flag.String("operation", "list", "It should accept there types of operation: [add|list|findById|remove], default: \"list\"")
	flag.Parse()
	return args
}

func Perform(args Arguments, writer io.Writer) error {
	var err error
	return err
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}
