package main

import (
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"

	pb "github.com/samertm/samerkv/samerkv"
	"google.golang.org/grpc"
)

var (
	port          = ":50051"
	storeFilename = filepath.Join(os.TempDir(), "default.skv")
	defaultTable  = "default"
)

func loadStore(path string) (map[string]kvTable, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Store doesn't exist on disk, return an
			// empty store.
			return map[string]kvTable{}, nil
		}
		return nil, err
	}
	defer f.Close()

	d := gob.NewDecoder(f)

	var store map[string]kvTable
	if err := d.Decode(&store); err != nil {
		return nil, err
	}

	return store, nil
}

func writeStore(path string, store map[string]kvTable) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	e := gob.NewEncoder(f)
	return e.Encode(store)
}

type kvTable map[string]string

var _ pb.UnstableKVStoreService = &grpcServer{}

type grpcServer struct {
	store map[string]kvTable
	mutex *sync.Mutex
}

func (g *grpcServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	table := req.GetTable()
	if table == "" {
		table = defaultTable
	}

	if g.store[table] == nil {
		return nil, fmt.Errorf("table doesn't exist: %s", table)
	}

	return &pb.GetResponse{Value: g.store[table][req.GetKey()]}, nil
}

func (g *grpcServer) Set(ctx context.Context, req *pb.SetRequest) (*pb.SetResponse, error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if req.GetKey() == "" {
		return nil, errors.New("key cannot be empty")
	}
	table := req.GetTable()
	if table == "" {
		table = defaultTable
	}

	// TODO: Figure out how to do this transactionally?
	if table == "*" {
		// Set value in all tables.
		for tableName := range g.store {
			g.store[tableName][req.GetKey()] = req.GetValue()
		}
	} else if g.store[table] == nil {
		return nil, fmt.Errorf("table doesn't exist: %s", table)
	} else {
		g.store[table][req.GetKey()] = req.GetValue()
	}

	if err := writeStore(storeFilename, g.store); err != nil {
		// TODO: rollback if there's an error
		return nil, fmt.Errorf("error writing store, database borked: %v", err)
	}
	return &pb.SetResponse{}, nil
}

func (g *grpcServer) CreateTable(ctx context.Context, req *pb.CreateTableRequest) (*pb.CreateTableResponse, error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if req.GetTable() == "" {
		return nil, errors.New("table name cannot be empty")
	}
	if g.store[req.GetTable()] != nil {
		return nil, fmt.Errorf("table named %s already exists", req.GetTable())
	}

	g.store[req.GetTable()] = kvTable{}
	if err := writeStore(storeFilename, g.store); err != nil {
		// TODO: rollback
		return nil, fmt.Errorf("error writing store, database borked: %v", err)
	}
	return &pb.CreateTableResponse{}, nil
}

func (g *grpcServer) DeleteTable(ctx context.Context, req *pb.DeleteTableRequest) (*pb.DeleteTableResponse, error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if req.GetTable() == "" {
		return nil, errors.New("table name cannot be empty")
	}
	if g.store[req.GetTable()] == nil {
		return nil, fmt.Errorf("table named %s does not exist", req.GetTable())
	}

	g.store[req.GetTable()] = nil
	if err := writeStore(storeFilename, g.store); err != nil {
		// TODO: rollback
		return nil, fmt.Errorf("error writing store, database borked: %v", err)
	}
	return &pb.DeleteTableResponse{}, nil
}

func (g *grpcServer) ListTables(ctx context.Context, req *pb.ListTablesRequest) (*pb.ListTablesResponse, error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	var tableNames []string
	for name := range g.store {
		tableNames = append(tableNames, name)
	}
	return &pb.ListTablesResponse{Tables: tableNames}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	store, err := loadStore(storeFilename)
	if err != nil {
		log.Fatalf("Error loading store: %v", err)
	}
	if len(store) == 0 {
		// Set up the default table
		store[defaultTable] = kvTable{}
	}

	server := &grpcServer{
		store: store,
		mutex: &sync.Mutex{},
	}

	s := grpc.NewServer()

	pb.RegisterKVStoreService(s, pb.NewKVStoreService(server))
	fmt.Println("Serving...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
