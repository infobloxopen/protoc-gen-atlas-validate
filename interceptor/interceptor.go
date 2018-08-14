package interceptor

import (
        "context"
	"fmt"
        "google.golang.org/grpc"
        "google.golang.org/grpc/codes"
        "google.golang.org/grpc/status"
        "google.golang.org/grpc/metadata"
)


const (
	ValidationErrorMetaKey = "Atlas-Validation-Error"
)

// PresenceClientInterceptor gets the interceptor for populating a fieldmask in a
// proto message from the fields given in the metadata/context
func ValidationClientInterceptor() grpc.UnaryClientInterceptor {
        return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
                if req == nil {
                        return
                }

		imd, _ := metadata.FromIncomingContext(ctx)
		omd, _ := metadata.FromOutgoingContext(ctx)

		md := metadata.Join(imd, omd)

		errors := md.Get(ValidationErrorMetaKey)
		if len(errors) > 0 {
			return status.Error(codes.InvalidArgument, errors[0])
		}

                return invoker(ctx, method, req, reply, cc, opts...)
        }
}

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
        return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (res interface{}, err error) {
		md, _ := metadata.FromIncomingContext(ctx)

		errors := md.Get(ValidationErrorMetaKey)
		if len(errors) > 0 {
			return nil, fmt.Errorf(errors[0])
		}

		return handler(ctx, req)
	}
}


