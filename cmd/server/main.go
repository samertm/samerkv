package main

import (
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"path/filepath"
	"sync"
	"unicode"
)

var (
	socketFilename = filepath.Join(os.TempDir(), "samerkv.sock")
	storeFilename  = filepath.Join(os.TempDir(), "default.skv")
)

type OperationType int

const (
	setOp OperationType = iota
	getOp OperationType = iota + 1
)

type Operation interface {
	IsOperation()
}

type getOperation struct {
	key string
}

func (getOperation) IsOperation() {}

type setOperation struct {
	key string
	val string
}

func (setOperation) IsOperation() {}

func lexInput(input string) ([]string, error) {
	var output []string
	var curWord string
	for _, c := range input {
		if unicode.IsSpace(c) || c == '=' {
			if len(curWord) != 0 {
				output = append(output, curWord)
				curWord = ""
			}
			if c == '=' {
				output = append(output, "=")
			}
		} else if unicode.IsLetter(c) || unicode.IsNumber(c) {
			curWord += string(c)
		} else {
			return nil, errors.New("Invalid character: " + string(c))
		}
	}
	if len(curWord) != 0 {
		output = append(output, curWord)
	}
	return output, nil
}

func parseInput(input string) (Operation, error) {
	// First, lex the input by turning the input stream into a
	// list of tokens.

	tokens, err := lexInput(input)
	if err != nil {
		return nil, err
	}
	if len(tokens) == 0 {
		return nil, nil
	}

	switch tokens[0] {
	case "get":
		if len(tokens) != 2 {
			return nil, errors.New("incorrect number of args for 'get'")
		}
		if tokens[1] == "=" {
			return nil, errors.New("cannot get '='")
		}
		return getOperation{key: tokens[1]}, nil
	case "set":
		if len(tokens) != 4 {
			return nil, errors.New("incorrect number of args for 'set'")
		}
		key := tokens[1]
		equal := tokens[2]
		val := tokens[3]
		if key == "=" {
			return nil, errors.New("key cannot be '='")
		}
		if equal != "=" {
			return nil, errors.New("expected '=', got: " + equal)
		}
		if val == "=" {
			return nil, errors.New("value cannot be '='")
		}
		return setOperation{key: key, val: val}, nil
	default:
		return nil, errors.New("invalid operation: " + tokens[0])
	}

	return nil, nil
}

func runOperation(store map[string]string, op Operation) (string, error) {
	switch o := op.(type) {
	case getOperation:
		return store[o.key], nil
	case setOperation:
		store[o.key] = o.val

		if err := writeStore(storeFilename, store); err != nil {
			return "", fmt.Errorf("Error writing store: %v", err)
		}
		return "", nil
	default:
		return "", fmt.Errorf("operation not implemented: %+v", op)
	}
}

func loadStore(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Store doesn't exist on disk, return an
			// empty store.
			return map[string]string{}, nil
		}
		return nil, err
	}
	defer f.Close()

	d := gob.NewDecoder(f)

	var store map[string]string
	if err := d.Decode(&store); err != nil {
		return nil, err
	}

	return store, nil
}

func writeStore(path string, store map[string]string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	e := gob.NewEncoder(f)
	return e.Encode(store)
}

type Server struct {
	store map[string]string
	mutex *sync.Mutex
}

func (s *Server) Execute(input *string, output *string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	operation, err := parseInput(*input)
	if err != nil {
		*output = fmt.Sprintf("error parsing input: %v", err)
		return nil
	}
	if operation == nil {
		*output = ""
		return nil
	}
	out, err := runOperation(s.store, operation)
	if err != nil {
		*output = fmt.Sprintf("error running operation: %v", err)
		return nil
	}
	*output = out
	return nil
}

func main() {
	if err := os.Remove(socketFilename); err != nil && !os.IsNotExist(err) {
		log.Fatalf("Could not remove socket %q: %v", socketFilename, err)
	}

	store, err := loadStore(storeFilename)
	if err != nil {
		log.Fatalf("Error loading store: %v", err)
	}

	l, err := net.Listen("unix", socketFilename)
	if err != nil {
		log.Fatalf("Listen error: %v", err)
	}
	server := Server{
		store: store,
		mutex: &sync.Mutex{},
	}

	rpc.Register(&server)
	rpc.HandleHTTP()
	fmt.Println("Serving...")
	http.Serve(l, nil)
}
