// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package examplepb

import (
	context "context"
	external "github.com/infobloxopen/protoc-gen-atlas-validate/example/external"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	anypb "google.golang.org/protobuf/types/known/anypb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// UsersClient is the client API for Users service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UsersClient interface {
	Create(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*EmptyResponse, error)
	Update(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (*EmptyResponse, error)
	List(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*EmptyResponse, error)
	UpdateExternalUser(ctx context.Context, in *User, opts ...grpc.CallOption) (*EmptyResponse, error)
	UpdateExternalUser2(ctx context.Context, in *external.ExternalUser, opts ...grpc.CallOption) (*EmptyResponse, error)
}

type usersClient struct {
	cc grpc.ClientConnInterface
}

func NewUsersClient(cc grpc.ClientConnInterface) UsersClient {
	return &usersClient{cc}
}

func (c *usersClient) Create(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*EmptyResponse, error) {
	out := new(EmptyResponse)
	err := c.cc.Invoke(ctx, "/examples.foo.Users/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersClient) Update(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (*EmptyResponse, error) {
	out := new(EmptyResponse)
	err := c.cc.Invoke(ctx, "/examples.foo.Users/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersClient) List(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*EmptyResponse, error) {
	out := new(EmptyResponse)
	err := c.cc.Invoke(ctx, "/examples.foo.Users/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersClient) UpdateExternalUser(ctx context.Context, in *User, opts ...grpc.CallOption) (*EmptyResponse, error) {
	out := new(EmptyResponse)
	err := c.cc.Invoke(ctx, "/examples.foo.Users/UpdateExternalUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersClient) UpdateExternalUser2(ctx context.Context, in *external.ExternalUser, opts ...grpc.CallOption) (*EmptyResponse, error) {
	out := new(EmptyResponse)
	err := c.cc.Invoke(ctx, "/examples.foo.Users/UpdateExternalUser2", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UsersServer is the server API for Users service.
// All implementations must embed UnimplementedUsersServer
// for forward compatibility
type UsersServer interface {
	Create(context.Context, *CreateUserRequest) (*EmptyResponse, error)
	Update(context.Context, *UpdateUserRequest) (*EmptyResponse, error)
	List(context.Context, *EmptyRequest) (*EmptyResponse, error)
	UpdateExternalUser(context.Context, *User) (*EmptyResponse, error)
	UpdateExternalUser2(context.Context, *external.ExternalUser) (*EmptyResponse, error)
	mustEmbedUnimplementedUsersServer()
}

// UnimplementedUsersServer must be embedded to have forward compatible implementations.
type UnimplementedUsersServer struct {
}

func (UnimplementedUsersServer) Create(context.Context, *CreateUserRequest) (*EmptyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedUsersServer) Update(context.Context, *UpdateUserRequest) (*EmptyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedUsersServer) List(context.Context, *EmptyRequest) (*EmptyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedUsersServer) UpdateExternalUser(context.Context, *User) (*EmptyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateExternalUser not implemented")
}
func (UnimplementedUsersServer) UpdateExternalUser2(context.Context, *external.ExternalUser) (*EmptyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateExternalUser2 not implemented")
}
func (UnimplementedUsersServer) mustEmbedUnimplementedUsersServer() {}

// UnsafeUsersServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UsersServer will
// result in compilation errors.
type UnsafeUsersServer interface {
	mustEmbedUnimplementedUsersServer()
}

func RegisterUsersServer(s grpc.ServiceRegistrar, srv UsersServer) {
	s.RegisterService(&Users_ServiceDesc, srv)
}

func _Users_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/examples.foo.Users/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServer).Create(ctx, req.(*CreateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Users_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/examples.foo.Users/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServer).Update(ctx, req.(*UpdateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Users_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/examples.foo.Users/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServer).List(ctx, req.(*EmptyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Users_UpdateExternalUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(User)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServer).UpdateExternalUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/examples.foo.Users/UpdateExternalUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServer).UpdateExternalUser(ctx, req.(*User))
	}
	return interceptor(ctx, in, info, handler)
}

func _Users_UpdateExternalUser2_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(external.ExternalUser)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServer).UpdateExternalUser2(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/examples.foo.Users/UpdateExternalUser2",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServer).UpdateExternalUser2(ctx, req.(*external.ExternalUser))
	}
	return interceptor(ctx, in, info, handler)
}

// Users_ServiceDesc is the grpc.ServiceDesc for Users service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Users_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "examples.foo.Users",
	HandlerType: (*UsersServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _Users_Create_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _Users_Update_Handler,
		},
		{
			MethodName: "List",
			Handler:    _Users_List_Handler,
		},
		{
			MethodName: "UpdateExternalUser",
			Handler:    _Users_UpdateExternalUser_Handler,
		},
		{
			MethodName: "UpdateExternalUser2",
			Handler:    _Users_UpdateExternalUser2_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "examplepb/example.proto",
}

// ProfilesClient is the client API for Profiles service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ProfilesClient interface {
	Create(ctx context.Context, in *Profile, opts ...grpc.CallOption) (*EmptyResponse, error)
	Update(ctx context.Context, in *UpdateProfileRequest, opts ...grpc.CallOption) (*EmptyResponse, error)
}

type profilesClient struct {
	cc grpc.ClientConnInterface
}

func NewProfilesClient(cc grpc.ClientConnInterface) ProfilesClient {
	return &profilesClient{cc}
}

func (c *profilesClient) Create(ctx context.Context, in *Profile, opts ...grpc.CallOption) (*EmptyResponse, error) {
	out := new(EmptyResponse)
	err := c.cc.Invoke(ctx, "/examples.foo.Profiles/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profilesClient) Update(ctx context.Context, in *UpdateProfileRequest, opts ...grpc.CallOption) (*EmptyResponse, error) {
	out := new(EmptyResponse)
	err := c.cc.Invoke(ctx, "/examples.foo.Profiles/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProfilesServer is the server API for Profiles service.
// All implementations must embed UnimplementedProfilesServer
// for forward compatibility
type ProfilesServer interface {
	Create(context.Context, *Profile) (*EmptyResponse, error)
	Update(context.Context, *UpdateProfileRequest) (*EmptyResponse, error)
	mustEmbedUnimplementedProfilesServer()
}

// UnimplementedProfilesServer must be embedded to have forward compatible implementations.
type UnimplementedProfilesServer struct {
}

func (UnimplementedProfilesServer) Create(context.Context, *Profile) (*EmptyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedProfilesServer) Update(context.Context, *UpdateProfileRequest) (*EmptyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedProfilesServer) mustEmbedUnimplementedProfilesServer() {}

// UnsafeProfilesServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ProfilesServer will
// result in compilation errors.
type UnsafeProfilesServer interface {
	mustEmbedUnimplementedProfilesServer()
}

func RegisterProfilesServer(s grpc.ServiceRegistrar, srv ProfilesServer) {
	s.RegisterService(&Profiles_ServiceDesc, srv)
}

func _Profiles_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Profile)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfilesServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/examples.foo.Profiles/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfilesServer).Create(ctx, req.(*Profile))
	}
	return interceptor(ctx, in, info, handler)
}

func _Profiles_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateProfileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfilesServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/examples.foo.Profiles/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfilesServer).Update(ctx, req.(*UpdateProfileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Profiles_ServiceDesc is the grpc.ServiceDesc for Profiles service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Profiles_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "examples.foo.Profiles",
	HandlerType: (*ProfilesServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _Profiles_Create_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _Profiles_Update_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "examplepb/example.proto",
}

// GroupsClient is the client API for Groups service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GroupsClient interface {
	Create(ctx context.Context, in *Group, opts ...grpc.CallOption) (*EmptyResponse, error)
	Update(ctx context.Context, in *Group, opts ...grpc.CallOption) (*EmptyResponse, error)
	ValidatedList(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*EmptyResponse, error)
	ValidateWKT(ctx context.Context, in *anypb.Any, opts ...grpc.CallOption) (*wrapperspb.DoubleValue, error)
}

type groupsClient struct {
	cc grpc.ClientConnInterface
}

func NewGroupsClient(cc grpc.ClientConnInterface) GroupsClient {
	return &groupsClient{cc}
}

func (c *groupsClient) Create(ctx context.Context, in *Group, opts ...grpc.CallOption) (*EmptyResponse, error) {
	out := new(EmptyResponse)
	err := c.cc.Invoke(ctx, "/examples.foo.Groups/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupsClient) Update(ctx context.Context, in *Group, opts ...grpc.CallOption) (*EmptyResponse, error) {
	out := new(EmptyResponse)
	err := c.cc.Invoke(ctx, "/examples.foo.Groups/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupsClient) ValidatedList(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*EmptyResponse, error) {
	out := new(EmptyResponse)
	err := c.cc.Invoke(ctx, "/examples.foo.Groups/ValidatedList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupsClient) ValidateWKT(ctx context.Context, in *anypb.Any, opts ...grpc.CallOption) (*wrapperspb.DoubleValue, error) {
	out := new(wrapperspb.DoubleValue)
	err := c.cc.Invoke(ctx, "/examples.foo.Groups/ValidateWKT", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GroupsServer is the server API for Groups service.
// All implementations must embed UnimplementedGroupsServer
// for forward compatibility
type GroupsServer interface {
	Create(context.Context, *Group) (*EmptyResponse, error)
	Update(context.Context, *Group) (*EmptyResponse, error)
	ValidatedList(context.Context, *EmptyRequest) (*EmptyResponse, error)
	ValidateWKT(context.Context, *anypb.Any) (*wrapperspb.DoubleValue, error)
	mustEmbedUnimplementedGroupsServer()
}

// UnimplementedGroupsServer must be embedded to have forward compatible implementations.
type UnimplementedGroupsServer struct {
}

func (UnimplementedGroupsServer) Create(context.Context, *Group) (*EmptyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedGroupsServer) Update(context.Context, *Group) (*EmptyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedGroupsServer) ValidatedList(context.Context, *EmptyRequest) (*EmptyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ValidatedList not implemented")
}
func (UnimplementedGroupsServer) ValidateWKT(context.Context, *anypb.Any) (*wrapperspb.DoubleValue, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ValidateWKT not implemented")
}
func (UnimplementedGroupsServer) mustEmbedUnimplementedGroupsServer() {}

// UnsafeGroupsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GroupsServer will
// result in compilation errors.
type UnsafeGroupsServer interface {
	mustEmbedUnimplementedGroupsServer()
}

func RegisterGroupsServer(s grpc.ServiceRegistrar, srv GroupsServer) {
	s.RegisterService(&Groups_ServiceDesc, srv)
}

func _Groups_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Group)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupsServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/examples.foo.Groups/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupsServer).Create(ctx, req.(*Group))
	}
	return interceptor(ctx, in, info, handler)
}

func _Groups_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Group)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupsServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/examples.foo.Groups/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupsServer).Update(ctx, req.(*Group))
	}
	return interceptor(ctx, in, info, handler)
}

func _Groups_ValidatedList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupsServer).ValidatedList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/examples.foo.Groups/ValidatedList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupsServer).ValidatedList(ctx, req.(*EmptyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Groups_ValidateWKT_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(anypb.Any)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupsServer).ValidateWKT(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/examples.foo.Groups/ValidateWKT",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupsServer).ValidateWKT(ctx, req.(*anypb.Any))
	}
	return interceptor(ctx, in, info, handler)
}

// Groups_ServiceDesc is the grpc.ServiceDesc for Groups service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Groups_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "examples.foo.Groups",
	HandlerType: (*GroupsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _Groups_Create_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _Groups_Update_Handler,
		},
		{
			MethodName: "ValidatedList",
			Handler:    _Groups_ValidatedList_Handler,
		},
		{
			MethodName: "ValidateWKT",
			Handler:    _Groups_ValidateWKT_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "examplepb/example.proto",
}
