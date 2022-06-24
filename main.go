package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type Arguments map[string]string

//type Arguments struct {
//operation string
//id        string
//item      string
//fileName  string
//}

type jsonuser struct {
	Id    string
	Email string
	Age   int
}

func parseArgs() Arguments {
	var args = Arguments{"operation": "", "id": "", "item": "", "fileName": ""}
	var operation, id, item, filename string
	//args["operation"] = *flag.String("operation", "", "It should accept there types of operation: [add|list|findById|remove], default: \"\"") //maybe default: list
	flag.StringVar(&operation, "operation", "", "It should accept there types of operation: [add|list|findById|remove].")
	//args["id"] = *flag.String("id", "", "user id, default: \"\"")
	flag.StringVar(&id, "id", "", "user id, default: \"\"")
	//args["item"] = *flag.String("item", "", "valid json object with the id, email and age fields") //maybe default: {\"id\": \"0\", \"email\": \"email@mail.com\", \"age\": 1}
	flag.StringVar(&item, "item", "", "valid json object with the id, email and age fields")
	//args["fileName"] = *flag.String("fileName", "", "name of file where you store users, default: \"users.json\"") //maybe default: users.json
	flag.StringVar(&filename, "fileName", "", "name of file where you store users, default: \"\"")
	flag.Parse()
	args["operation"] = operation
	args["id"] = id
	args["item"] = item
	args["fileName"] = filename
	return args
}

func AddingOperation(jsonstring string, filename string, wrt io.Writer) error {
	var users []jsonuser
	var user jsonuser
	err := json.Unmarshal([]byte(jsonstring), &user)
	if err != nil {
		return fmt.Errorf("json format incorrect: %w", err)
	}
	content, _ := os.ReadFile(filename)
	err = json.Unmarshal(content, &users)
	if err == nil && len(users) != 0 {
		for _, usr := range users {
			if user.Id == usr.Id {
				wrt.Write([]byte("Item with id " + usr.Id + " already exists"))
				return nil
			}
		}
	}
	users = append(users, user)
	content, err = json.Marshal(users)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return fmt.Errorf("%w", err)
	} else {
		_, err = file.Write(content)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
		err = file.Close()
		if err != nil {
			return fmt.Errorf("%w", err)
		}
	}
	return nil
}

func ListOperation(filename string, wrt io.Writer) error {
	var users []jsonuser
	file, err := os.OpenFile(filename, os.O_RDONLY, 0755)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	err = json.Unmarshal(content, &users)
	if err != nil {
		return fmt.Errorf("%w", err)
	} else {
		wrt.Write(content)
		return nil
	}
	//return nil
}

func FindByIdOperation(filename string, id string, wrt io.Writer) error {
	var users []jsonuser
	var temp []byte
	file, err := os.OpenFile(filename, os.O_RDONLY, 0755)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	err = json.Unmarshal(content, &users)
	if err == nil && len(users) != 0 {
		for _, usr := range users {
			if id == usr.Id {
				temp, err = json.Marshal(usr)
				if err != nil {
					return fmt.Errorf("%w", err)
				}
				wrt.Write(temp)
				return nil
			}
		}
	}
	return nil
}

func RemovingOperation(filename string, id string, wrt io.Writer) error {
	var users []jsonuser
	//var temp []byte
	file, err := os.OpenFile(filename, os.O_RDONLY, 0755)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	err = json.Unmarshal(content, &users)
	if err == nil && len(users) != 0 {
		found := false
		for i, usr := range users {
			if id == usr.Id {
				found = true
				//temp, err = json.Marshal(usr)
				//if err != nil {
				//return fmt.Errorf("%w", err)
				//}
				//wrt.Write(temp)
				users = append(users[:i], users[i+1:]...)
				content, err = json.Marshal(users)
				if err != nil {
					return fmt.Errorf("%w", err)
				}
				file, err = os.OpenFile(filename, os.O_WRONLY, 0755)
				if err != nil {
					return fmt.Errorf("%w", err)
				}
				_, err = file.Write(content)
				if err != nil {
					return fmt.Errorf("%w", err)
				}
				err = file.Close()
				if err != nil {
					return fmt.Errorf("%w", err)
				}
			}
		}
		if !found {
			return fmt.Errorf("Item with id " + id + " not found")
		}
	}
	return nil
}

func Perform(args Arguments, writer io.Writer) error {
	//var err error
	if args["fileName"] == "" {
		return errors.New("-fileName flag has to be specified")
	}
	if args["operation"] == "" {
		return errors.New("-operation flag has to be specified")
	}
	switch args["operation"] {
	case "add":
		if args["item"] == "" {
			return errors.New("-item flag has to be specified")
		} else {
			return AddingOperation(args["item"], args["fileName"], writer)
		}
	case "list":
		return ListOperation(args["fileName"], writer)
	case "findById":
		if args["id"] == "" {
			return errors.New("-id flag has to be specified")
		} else {
			return nil
		}
	case "remove":
		if args["id"] == "" {
			return errors.New("-id flag has to be specified")
		} else {
			return nil
		}
	default:
		return errors.New("Operation " + args["operation"] + " not allowed!")
	}
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}
