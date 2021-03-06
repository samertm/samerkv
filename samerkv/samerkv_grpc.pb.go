// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package samerkv

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// KVStoreClient is the client API for KVStore service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type KVStoreClient interface {
	Set(ctx context.Context, in *SetRequest, opts ...grpc.CallOption) (*SetResponse, error)
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error)
	CreateTable(ctx context.Context, in *CreateTableRequest, opts ...grpc.CallOption) (*CreateTableResponse, error)
	DeleteTable(ctx context.Context, in *DeleteTableRequest, opts ...grpc.CallOption) (*DeleteTableResponse, error)
	ListTables(ctx context.Context, in *ListTablesRequest, opts ...grpc.CallOption) (*ListTablesResponse, error)
}

type kVStoreClient struct {
	cc grpc.ClientConnInterface
}

func NewKVStoreClient(cc grpc.ClientConnInterface) KVStoreClient {
	return &kVStoreClient{cc}
}

var kVStoreSetStreamDesc = &grpc.StreamDesc{
	StreamName: "Set",
}

func (c *kVStoreClient) Set(ctx context.Context, in *SetRequest, opts ...grpc.CallOption) (*SetResponse, error) {
	out := new(SetResponse)
	err := c.cc.Invoke(ctx, "/samerkv.KVStore/Set", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var kVStoreGetStreamDesc = &grpc.StreamDesc{
	StreamName: "Get",
}

func (c *kVStoreClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	out := new(GetResponse)
	err := c.cc.Invoke(ctx, "/samerkv.KVStore/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var kVStoreCreateTableStreamDesc = &grpc.StreamDesc{
	StreamName: "CreateTable",
}

func (c *kVStoreClient) CreateTable(ctx context.Context, in *CreateTableRequest, opts ...grpc.CallOption) (*CreateTableResponse, error) {
	out := new(CreateTableResponse)
	err := c.cc.Invoke(ctx, "/samerkv.KVStore/CreateTable", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var kVStoreDeleteTableStreamDesc = &grpc.StreamDesc{
	StreamName: "DeleteTable",
}

func (c *kVStoreClient) DeleteTable(ctx context.Context, in *DeleteTableRequest, opts ...grpc.CallOption) (*DeleteTableResponse, error) {
	out := new(DeleteTableResponse)
	err := c.cc.Invoke(ctx, "/samerkv.KVStore/DeleteTable", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var kVStoreListTablesStreamDesc = &grpc.StreamDesc{
	StreamName: "ListTables",
}

func (c *kVStoreClient) ListTables(ctx context.Context, in *ListTablesRequest, opts ...grpc.CallOption) (*ListTablesResponse, error) {
	out := new(ListTablesResponse)
	err := c.cc.Invoke(ctx, "/samerkv.KVStore/ListTables", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// KVStoreService is the service API for KVStore service.
// Fields should be assigned to their respective handler implementations only before
// RegisterKVStoreService is called.  Any unassigned fields will result in the
// handler for that method returning an Unimplemented error.
type KVStoreService struct {
	Set         func(context.Context, *SetRequest) (*SetResponse, error)
	Get         func(context.Context, *GetRequest) (*GetResponse, error)
	CreateTable func(context.Context, *CreateTableRequest) (*CreateTableResponse, error)
	DeleteTable func(context.Context, *DeleteTableRequest) (*DeleteTableResponse, error)
	ListTables  func(context.Context, *ListTablesRequest) (*ListTablesResponse, error)
}

func (s *KVStoreService) set(_ interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return s.Set(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     s,
		FullMethod: "/samerkv.KVStore/Set",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return s.Set(ctx, req.(*SetRequest))
	}
	return interceptor(ctx, in, info, handler)
}
func (s *KVStoreService) get(_ interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return s.Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     s,
		FullMethod: "/samerkv.KVStore/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return s.Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}
func (s *KVStoreService) createTable(_ interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateTableRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return s.CreateTable(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     s,
		FullMethod: "/samerkv.KVStore/CreateTable",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return s.CreateTable(ctx, req.(*CreateTableRequest))
	}
	return interceptor(ctx, in, info, handler)
}
func (s *KVStoreService) deleteTable(_ interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteTableRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return s.DeleteTable(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     s,
		FullMethod: "/samerkv.KVStore/DeleteTable",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return s.DeleteTable(ctx, req.(*DeleteTableRequest))
	}
	return interceptor(ctx, in, info, handler)
}
func (s *KVStoreService) listTables(_ interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListTablesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return s.ListTables(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     s,
		FullMethod: "/samerkv.KVStore/ListTables",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return s.ListTables(ctx, req.(*ListTablesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RegisterKVStoreService registers a service implementation with a gRPC server.
func RegisterKVStoreService(s grpc.ServiceRegistrar, srv *KVStoreService) {
	srvCopy := *srv
	if srvCopy.Set == nil {
		srvCopy.Set = func(context.Context, *SetRequest) (*SetResponse, error) {
			return nil, status.Errorf(codes.Unimplemented, "method Set not implemented")
		}
	}
	if srvCopy.Get == nil {
		srvCopy.Get = func(context.Context, *GetRequest) (*GetResponse, error) {
			return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
		}
	}
	if srvCopy.CreateTable == nil {
		srvCopy.CreateTable = func(context.Context, *CreateTableRequest) (*CreateTableResponse, error) {
			return nil, status.Errorf(codes.Unimplemented, "method CreateTable not implemented")
		}
	}
	if srvCopy.DeleteTable == nil {
		srvCopy.DeleteTable = func(context.Context, *DeleteTableRequest) (*DeleteTableResponse, error) {
			return nil, status.Errorf(codes.Unimplemented, "method DeleteTable not implemented")
		}
	}
	if srvCopy.ListTables == nil {
		srvCopy.ListTables = func(context.Context, *ListTablesRequest) (*ListTablesResponse, error) {
			return nil, status.Errorf(codes.Unimplemented, "method ListTables not implemented")
		}
	}
	sd := grpc.ServiceDesc{
		ServiceName: "samerkv.KVStore",
		Methods: []grpc.MethodDesc{
			{
				MethodName: "Set",
				Handler:    srvCopy.set,
			},
			{
				MethodName: "Get",
				Handler:    srvCopy.get,
			},
			{
				MethodName: "CreateTable",
				Handler:    srvCopy.createTable,
			},
			{
				MethodName: "DeleteTable",
				Handler:    srvCopy.deleteTable,
			},
			{
				MethodName: "ListTables",
				Handler:    srvCopy.listTables,
			},
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "samerkv/samerkv.proto",
	}

	s.RegisterService(&sd, nil)
}

// NewKVStoreService creates a new KVStoreService containing the
// implemented methods of the KVStore service in s.  Any unimplemented
// methods will result in the gRPC server returning an UNIMPLEMENTED status to the client.
// This includes situations where the method handler is misspelled or has the wrong
// signature.  For this reason, this function should be used with great care and
// is not recommended to be used by most users.
func NewKVStoreService(s interface{}) *KVStoreService {
	ns := &KVStoreService{}
	if h, ok := s.(interface {
		Set(context.Context, *SetRequest) (*SetResponse, error)
	}); ok {
		ns.Set = h.Set
	}
	if h, ok := s.(interface {
		Get(context.Context, *GetRequest) (*GetResponse, error)
	}); ok {
		ns.Get = h.Get
	}
	if h, ok := s.(interface {
		CreateTable(context.Context, *CreateTableRequest) (*CreateTableResponse, error)
	}); ok {
		ns.CreateTable = h.CreateTable
	}
	if h, ok := s.(interface {
		DeleteTable(context.Context, *DeleteTableRequest) (*DeleteTableResponse, error)
	}); ok {
		ns.DeleteTable = h.DeleteTable
	}
	if h, ok := s.(interface {
		ListTables(context.Context, *ListTablesRequest) (*ListTablesResponse, error)
	}); ok {
		ns.ListTables = h.ListTables
	}
	return ns
}

// UnstableKVStoreService is the service API for KVStore service.
// New methods may be added to this interface if they are added to the service
// definition, which is not a backward-compatible change.  For this reason,
// use of this type is not recommended.
type UnstableKVStoreService interface {
	Set(context.Context, *SetRequest) (*SetResponse, error)
	Get(context.Context, *GetRequest) (*GetResponse, error)
	CreateTable(context.Context, *CreateTableRequest) (*CreateTableResponse, error)
	DeleteTable(context.Context, *DeleteTableRequest) (*DeleteTableResponse, error)
	ListTables(context.Context, *ListTablesRequest) (*ListTablesResponse, error)
}

// InternalKVStoreClient is the client API for InternalKVStore service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type InternalKVStoreClient interface {
	GetLogCount(ctx context.Context, in *GetLogCountRequest, opts ...grpc.CallOption) (*GetLogCountResponse, error)
	ReplicateLog(ctx context.Context, in *ReplicateLogRequest, opts ...grpc.CallOption) (*ReplicateLogResponse, error)
}

type internalKVStoreClient struct {
	cc grpc.ClientConnInterface
}

func NewInternalKVStoreClient(cc grpc.ClientConnInterface) InternalKVStoreClient {
	return &internalKVStoreClient{cc}
}

var internalKVStoreGetLogCountStreamDesc = &grpc.StreamDesc{
	StreamName: "GetLogCount",
}

func (c *internalKVStoreClient) GetLogCount(ctx context.Context, in *GetLogCountRequest, opts ...grpc.CallOption) (*GetLogCountResponse, error) {
	out := new(GetLogCountResponse)
	err := c.cc.Invoke(ctx, "/samerkv.InternalKVStore/GetLogCount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var internalKVStoreReplicateLogStreamDesc = &grpc.StreamDesc{
	StreamName: "ReplicateLog",
}

func (c *internalKVStoreClient) ReplicateLog(ctx context.Context, in *ReplicateLogRequest, opts ...grpc.CallOption) (*ReplicateLogResponse, error) {
	out := new(ReplicateLogResponse)
	err := c.cc.Invoke(ctx, "/samerkv.InternalKVStore/ReplicateLog", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// InternalKVStoreService is the service API for InternalKVStore service.
// Fields should be assigned to their respective handler implementations only before
// RegisterInternalKVStoreService is called.  Any unassigned fields will result in the
// handler for that method returning an Unimplemented error.
type InternalKVStoreService struct {
	GetLogCount  func(context.Context, *GetLogCountRequest) (*GetLogCountResponse, error)
	ReplicateLog func(context.Context, *ReplicateLogRequest) (*ReplicateLogResponse, error)
}

func (s *InternalKVStoreService) getLogCount(_ interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetLogCountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return s.GetLogCount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     s,
		FullMethod: "/samerkv.InternalKVStore/GetLogCount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return s.GetLogCount(ctx, req.(*GetLogCountRequest))
	}
	return interceptor(ctx, in, info, handler)
}
func (s *InternalKVStoreService) replicateLog(_ interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReplicateLogRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return s.ReplicateLog(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     s,
		FullMethod: "/samerkv.InternalKVStore/ReplicateLog",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return s.ReplicateLog(ctx, req.(*ReplicateLogRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RegisterInternalKVStoreService registers a service implementation with a gRPC server.
func RegisterInternalKVStoreService(s grpc.ServiceRegistrar, srv *InternalKVStoreService) {
	srvCopy := *srv
	if srvCopy.GetLogCount == nil {
		srvCopy.GetLogCount = func(context.Context, *GetLogCountRequest) (*GetLogCountResponse, error) {
			return nil, status.Errorf(codes.Unimplemented, "method GetLogCount not implemented")
		}
	}
	if srvCopy.ReplicateLog == nil {
		srvCopy.ReplicateLog = func(context.Context, *ReplicateLogRequest) (*ReplicateLogResponse, error) {
			return nil, status.Errorf(codes.Unimplemented, "method ReplicateLog not implemented")
		}
	}
	sd := grpc.ServiceDesc{
		ServiceName: "samerkv.InternalKVStore",
		Methods: []grpc.MethodDesc{
			{
				MethodName: "GetLogCount",
				Handler:    srvCopy.getLogCount,
			},
			{
				MethodName: "ReplicateLog",
				Handler:    srvCopy.replicateLog,
			},
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "samerkv/samerkv.proto",
	}

	s.RegisterService(&sd, nil)
}

// NewInternalKVStoreService creates a new InternalKVStoreService containing the
// implemented methods of the InternalKVStore service in s.  Any unimplemented
// methods will result in the gRPC server returning an UNIMPLEMENTED status to the client.
// This includes situations where the method handler is misspelled or has the wrong
// signature.  For this reason, this function should be used with great care and
// is not recommended to be used by most users.
func NewInternalKVStoreService(s interface{}) *InternalKVStoreService {
	ns := &InternalKVStoreService{}
	if h, ok := s.(interface {
		GetLogCount(context.Context, *GetLogCountRequest) (*GetLogCountResponse, error)
	}); ok {
		ns.GetLogCount = h.GetLogCount
	}
	if h, ok := s.(interface {
		ReplicateLog(context.Context, *ReplicateLogRequest) (*ReplicateLogResponse, error)
	}); ok {
		ns.ReplicateLog = h.ReplicateLog
	}
	return ns
}

// UnstableInternalKVStoreService is the service API for InternalKVStore service.
// New methods may be added to this interface if they are added to the service
// definition, which is not a backward-compatible change.  For this reason,
// use of this type is not recommended.
type UnstableInternalKVStoreService interface {
	GetLogCount(context.Context, *GetLogCountRequest) (*GetLogCountResponse, error)
	ReplicateLog(context.Context, *ReplicateLogRequest) (*ReplicateLogResponse, error)
}
