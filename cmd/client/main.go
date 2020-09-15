package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/peterh/liner"
	pb "github.com/samertm/samerkv/samerkv"
	"google.golang.org/grpc"
)

var (
	address         = "localhost:50051"
	historyFilename = filepath.Join(os.TempDir(), "samerkv_history")
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
	key   string
	table string
}

func (getOperation) IsOperation() {}

type setOperation struct {
	key   string
	table string
	val   string
}

func (setOperation) IsOperation() {}

type createTableOperation struct {
	name string
}

func (createTableOperation) IsOperation() {}

type deleteTableOperation struct {
	name string
}

func (deleteTableOperation) IsOperation() {}

type listTablesOperation struct{}

func (listTablesOperation) IsOperation() {}

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
		} else if unicode.IsLetter(c) || unicode.IsNumber(c) || c == '.' || c == '*' || c == '_' {
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

func parseKeyAndTable(rawKey string) (key string, table string, err error) {
	splitKey := strings.Split(rawKey, ".")
	if len(splitKey) > 2 {
		return "", "", errors.New("too many '.'")
	} else if len(splitKey) == 2 {
		table = splitKey[0]
		key = splitKey[1]
	} else {
		key = splitKey[0]
	}
	return key, table, nil
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
		key, table, err := parseKeyAndTable(tokens[1])
		if err != nil {
			return nil, err
		}

		return getOperation{key: key, table: table}, nil
	case "set":
		if len(tokens) != 4 {
			return nil, errors.New("incorrect number of args for 'set'")
		}
		rawKey := tokens[1]
		equal := tokens[2]
		val := tokens[3]
		if rawKey == "=" {
			return nil, errors.New("key cannot be '='")
		}
		if equal != "=" {
			return nil, errors.New("expected '=', got: " + equal)
		}
		if val == "=" {
			return nil, errors.New("value cannot be '='")
		}
		key, table, err := parseKeyAndTable(rawKey)
		if err != nil {
			return nil, err
		}
		return setOperation{key: key, table: table, val: val}, nil
	case "create_table":
		if len(tokens) != 2 {
			return nil, errors.New("incorrect number of args for 'create_table'")
		}
		if tokens[1] == "=" {
			return nil, errors.New("table name cannot be '='")
		}
		return createTableOperation{name: tokens[1]}, nil
	case "delete_table":
		if len(tokens) != 2 {
			return nil, errors.New("incorrect number of args for 'delete_table'")
		}
		if tokens[1] == "=" {
			return nil, errors.New("table name cannot be '='")
		}
		return deleteTableOperation{name: tokens[1]}, nil
	case "list_tables":
		if len(tokens) != 1 {
			return nil, errors.New("incorrect number of args for 'list_table'")
		}
		return listTablesOperation{}, nil
	default:
		return nil, errors.New("invalid operation: " + tokens[0])
	}

	return nil, nil
}

func runOperation(client pb.KVStoreClient, op Operation) (string, error) {
	switch o := op.(type) {
	case getOperation:
		resp, err := client.Get(context.Background(), &pb.GetRequest{
			Key:   o.key,
			Table: o.table,
		})
		if err != nil {
			return "", err
		}
		return resp.GetValue(), nil
	case setOperation:
		_, err := client.Set(context.Background(), &pb.SetRequest{
			Key:   o.key,
			Value: o.val,
			Table: o.table,
		})
		if err != nil {
			return "", err
		}
		return "", nil
	case createTableOperation:
		_, err := client.CreateTable(context.Background(), &pb.CreateTableRequest{
			Table: o.name,
		})
		if err != nil {
			return "", err
		}
		return "", nil
	case deleteTableOperation:
		_, err := client.DeleteTable(context.Background(), &pb.DeleteTableRequest{
			Table: o.name,
		})
		if err != nil {
			return "", err
		}
		return "", nil
	case listTablesOperation:
		resp, err := client.ListTables(context.Background(), &pb.ListTablesRequest{})
		if err != nil {
			return "", err
		}
		var output string
		for _, table := range resp.GetTables() {
			if output != "" {
				output += ", "
			}
			output += table
		}
		return output, err
	default:
		return "", fmt.Errorf("operation not implemented: %+v", op)
	}
}

func main() {
	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)

	if f, err := os.Open(historyFilename); err == nil {
		line.ReadHistory(f)
		f.Close()
	}

	// Establish connection.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewKVStoreClient(conn)

	for {
		input, err := line.Prompt("> ")
		if err != nil {
			if err == liner.ErrPromptAborted || err == io.EOF {
				break
			}
			log.Printf("error reading prompt: %v", err)
			continue
		}
		line.AppendHistory(input)

		op, err := parseInput(input)
		if err != nil {
			log.Printf("error parsing input: %v", err)
			continue
		}
		if op == nil {
			continue
		}
		output, err := runOperation(client, op)
		if err != nil {
			log.Printf("error running operation: %v", err)
			continue
		}
		if len(output) > 0 {
			fmt.Println(output)
		}
	}

	if f, err := os.Create(historyFilename); err != nil {
		log.Printf("Error writing history file: %v", err)
	} else {
		line.WriteHistory(f)
		f.Close()
	}
}
