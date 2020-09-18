package main

import (
	"context"
	"encoding/gob"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"

	pb "github.com/samertm/samerkv/samerkv"
	"google.golang.org/grpc"
)

const (
	KindLeader        = "leader"
	KindSyncFollower  = "syncfollower"
	KindAsyncFollower = "asyncfollower"
)

var (
	asyncFollowerPort = ":50053"
	syncFollowerPort  = ":50052"
	leaderPort        = ":50051"
	defaultTable      = "default"

	flagKind string
)

func init() {
	flag.StringVar(&flagKind, "kind", "", "Available values are 'leader', 'syncfollower', and 'asyncfollower'. Required.")
}

type kvTable map[string]string

type kvStore struct {
	Data map[string]kvTable
	Logs []*pb.WrappedRequest
}

func loadStore(path string) (*kvStore, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Store doesn't exist on disk, return an
			// empty store.
			return &kvStore{Data: map[string]kvTable{}}, nil
		}
		return nil, err
	}
	defer f.Close()

	d := gob.NewDecoder(f)

	var store kvStore
	if err := d.Decode(&store); err != nil {
		return nil, err
	}

	return &store, nil
}

func writeStore(path string, store *kvStore) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	e := gob.NewEncoder(f)
	return e.Encode(*store)
}

var _ pb.UnstableKVStoreService = &grpcServer{}

type grpcServer struct {
	// leader is true if this server is the leader.
	leader bool
	// syncFollowerClient is only set if this server is the leader.
	syncFollowerClient pb.InternalKVStoreClient
	// asyncFollowerClient is only set if this server is the leader.
	asyncFollowerClient pb.InternalKVStoreClient
	storeFilename       string
	store               *kvStore
	mutex               *sync.Mutex
}

func (g *grpcServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	table := req.GetTable()
	if table == "" {
		table = defaultTable
	}

	if g.store.Data[table] == nil {
		return nil, fmt.Errorf("table doesn't exist: %s", table)
	}

	return &pb.GetResponse{Value: g.store.Data[table][req.GetKey()]}, nil
}

func (g *grpcServer) Set(ctx context.Context, req *pb.SetRequest) (*pb.SetResponse, error) {
	if !g.leader {
		return nil, errors.New("cannot make write request to follower")
	}
	g.mutex.Lock()
	defer g.mutex.Unlock()
	return g.internalSet(ctx, req)
}

func (g *grpcServer) internalSet(ctx context.Context, req *pb.SetRequest) (*pb.SetResponse, error) {
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
		for tableName := range g.store.Data {
			g.store.Data[tableName][req.GetKey()] = req.GetValue()
		}
	} else if g.store.Data[table] == nil {
		return nil, fmt.Errorf("table doesn't exist: %s", table)
	} else {
		g.store.Data[table][req.GetKey()] = req.GetValue()
	}
	g.store.Logs = append(g.store.Logs, &pb.WrappedRequest{
		Req: &pb.WrappedRequest_SetReq{SetReq: req}})

	if err := writeStore(g.storeFilename, g.store); err != nil {
		// TODO: rollback if there's an error
		return nil, fmt.Errorf("error writing store, database borked: %v", err)
	}
	if g.leader {
		if err := sendLogsToFollowers(ctx, g); err != nil {
			return nil, err
		}
	}
	return &pb.SetResponse{}, nil
}

func (g *grpcServer) CreateTable(ctx context.Context, req *pb.CreateTableRequest) (*pb.CreateTableResponse, error) {
	if !g.leader {
		return nil, errors.New("cannot make write request to follower")
	}
	g.mutex.Lock()
	defer g.mutex.Unlock()
	return g.internalCreateTable(ctx, req)
}

func (g *grpcServer) internalCreateTable(ctx context.Context, req *pb.CreateTableRequest) (*pb.CreateTableResponse, error) {
	if req.GetTable() == "" {
		return nil, errors.New("table name cannot be empty")
	}
	if g.store.Data[req.GetTable()] != nil {
		return nil, fmt.Errorf("table named %s already exists", req.GetTable())
	}

	g.store.Data[req.GetTable()] = kvTable{}
	g.store.Logs = append(g.store.Logs, &pb.WrappedRequest{Req: &pb.WrappedRequest_CreateTableReq{CreateTableReq: req}})
	if err := writeStore(g.storeFilename, g.store); err != nil {
		// TODO: rollback
		return nil, fmt.Errorf("error writing store, database borked: %v", err)
	}
	if g.leader {
		if err := sendLogsToFollowers(ctx, g); err != nil {
			return nil, err
		}
	}
	return &pb.CreateTableResponse{}, nil
}

func (g *grpcServer) DeleteTable(ctx context.Context, req *pb.DeleteTableRequest) (*pb.DeleteTableResponse, error) {
	if !g.leader {
		return nil, errors.New("cannot make write request to follower")
	}
	g.mutex.Lock()
	defer g.mutex.Unlock()
	return g.internalDeleteTable(ctx, req)
}

func (g *grpcServer) internalDeleteTable(ctx context.Context, req *pb.DeleteTableRequest) (*pb.DeleteTableResponse, error) {
	if req.GetTable() == "" {
		return nil, errors.New("table name cannot be empty")
	}
	if g.store.Data[req.GetTable()] == nil {
		return nil, fmt.Errorf("table named %s does not exist", req.GetTable())
	}

	g.store.Data[req.GetTable()] = nil
	g.store.Logs = append(g.store.Logs, &pb.WrappedRequest{
		Req: &pb.WrappedRequest_DeleteTableReq{DeleteTableReq: req}})
	if err := writeStore(g.storeFilename, g.store); err != nil {
		// TODO: rollback
		return nil, fmt.Errorf("error writing store, database borked: %v", err)
	}
	if g.leader {
		if err := sendLogsToFollowers(ctx, g); err != nil {
			return nil, err
		}
	}
	return &pb.DeleteTableResponse{}, nil
}

func (g *grpcServer) ListTables(ctx context.Context, req *pb.ListTablesRequest) (*pb.ListTablesResponse, error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	var tableNames []string
	for name := range g.store.Data {
		tableNames = append(tableNames, name)
	}
	return &pb.ListTablesResponse{Tables: tableNames}, nil
}

func (g *grpcServer) GetLogCount(ctx context.Context, req *pb.GetLogCountRequest) (*pb.GetLogCountResponse, error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	return &pb.GetLogCountResponse{LogCount: int32(len(g.store.Logs))}, nil
}

func (g *grpcServer) ReplicateLog(ctx context.Context, req *pb.ReplicateLogRequest) (*pb.ReplicateLogResponse, error) {
	if g.leader {
		return nil, errors.New("Cannot make ReplicateLog request to leader")
	}
	g.mutex.Lock()
	defer g.mutex.Unlock()

	startCount := int32(len(g.store.Logs))
	if startCount > req.StartAt {
		return nil, fmt.Errorf("invalid start at. StartAt: %d, start count: %d", req.StartAt, startCount)
	}

	for i, logProto := range req.Requests {
		if req.StartAt+int32(i) < startCount {
			continue
		}
		switch p := logProto.Req.(type) {
		case *pb.WrappedRequest_SetReq:
			_, err := g.internalSet(ctx, p.SetReq)
			if err != nil {
				return nil, err
			}
		case *pb.WrappedRequest_CreateTableReq:
			_, err := g.internalCreateTable(ctx, p.CreateTableReq)
			if err != nil {
				return nil, err
			}
		case *pb.WrappedRequest_DeleteTableReq:
			_, err := g.internalDeleteTable(ctx, p.DeleteTableReq)
			if err != nil {
				return nil, err
			}
		default:
			return nil, errors.New("unexpected proto, follower in invalid state")
		}
	}
	return &pb.ReplicateLogResponse{}, nil
}

func sendLogsToFollowers(ctx context.Context, g *grpcServer) error {
	log.Print("sending logs to sync follower")
	logs := make([]*pb.WrappedRequest, len(g.store.Logs))
	copy(logs, g.store.Logs)
	if err := sendLogsToFollower(ctx, g.syncFollowerClient, logs); err != nil {
		return err
	}
	// Make the second request async
	go func() {
		if err := sendLogsToFollower(context.TODO(), g.asyncFollowerClient, logs); err != nil {
			log.Printf("Error writing some logs to async client: %+v, %v", logs, err)
		}
	}()
	return nil
}

func sendLogsToFollower(ctx context.Context, client pb.InternalKVStoreClient, logs []*pb.WrappedRequest) error {
	resp, err := client.GetLogCount(ctx, &pb.GetLogCountRequest{})
	if err != nil {
		return err
	}
	currentLogCount := int32(len(logs))
	if resp.LogCount > currentLogCount {
		return fmt.Errorf("Follow log count greater than leader. follower: %d, leader: %d",
			resp.LogCount, len(logs))
	}
	if resp.LogCount == currentLogCount {
		return nil
	}

	// Send missing logs
	var toSend []*pb.WrappedRequest
	for i := resp.LogCount; i < currentLogCount; i++ {
		toSend = append(toSend, logs[i])
	}
	_, err = client.ReplicateLog(ctx, &pb.ReplicateLogRequest{
		StartAt:  resp.LogCount,
		Requests: toSend,
	})
	return err
}

func internalClient(port string) pb.InternalKVStoreClient {
	conn, err := grpc.Dial(port, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect to follower: %v", err)
	}
	return pb.NewInternalKVStoreClient(conn)
}

func main() {
	gob.Register(&pb.WrappedRequest_SetReq{})
	gob.Register(&pb.WrappedRequest_CreateTableReq{})
	gob.Register(&pb.WrappedRequest_DeleteTableReq{})
	flag.Parse()
	if flagKind != KindLeader && flagKind != KindSyncFollower && flagKind != KindAsyncFollower {
		log.Fatalf("Invalid value for --kind: %s", flagKind)
	}
	var port string
	var storeFilename string
	switch flagKind {
	case KindLeader:
		fmt.Println("running leader")
		port = leaderPort
		storeFilename = filepath.Join(os.TempDir(), "default.skv")
	case KindSyncFollower:
		fmt.Println("running sync follower")
		port = syncFollowerPort
		storeFilename = filepath.Join(os.TempDir(), "default-syncfollower.skv")
	case KindAsyncFollower:
		fmt.Println("running async follower")
		port = asyncFollowerPort
		storeFilename = filepath.Join(os.TempDir(), "default-asyncfollower.skv")
	default:
		log.Fatalf("unreachable")
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	store, err := loadStore(storeFilename)
	if err != nil {
		log.Fatalf("Error loading store: %v", err)
	}
	if len(store.Logs) == 0 {
		// Set up the default table
		store.Data[defaultTable] = kvTable{}
	}
	var syncFollowerClient pb.InternalKVStoreClient
	var asyncFollowerClient pb.InternalKVStoreClient
	if flagKind == KindLeader {
		syncFollowerClient = internalClient(syncFollowerPort)
		asyncFollowerClient = internalClient(asyncFollowerPort)
	}

	server := &grpcServer{
		leader:              flagKind == KindLeader,
		syncFollowerClient:  syncFollowerClient,
		asyncFollowerClient: asyncFollowerClient,
		storeFilename:       storeFilename,
		store:               store,
		mutex:               &sync.Mutex{},
	}

	s := grpc.NewServer()

	pb.RegisterKVStoreService(s, pb.NewKVStoreService(server))
	pb.RegisterInternalKVStoreService(s, pb.NewInternalKVStoreService(server))

	if flagKind == KindLeader {
		if err := sendLogsToFollowers(context.Background(), server); err != nil {
			log.Fatalf("error replicating follower: %v", err)
		}
	}
	fmt.Println("Serving...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
