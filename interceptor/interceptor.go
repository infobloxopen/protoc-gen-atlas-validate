package interceptor

import (
        "context"
	"fmt"
        "google.golang.org/grpc"
        "google.golang.org/grpc/metadata"
)


const (
	ValidationErrorMetaKey = "Validation-Error"
)

// PresenceClientInterceptor gets the interceptor for populating a fieldmask in a
// proto message from the fields given in the metadata/context
func ValidationClientInterceptor() grpc.UnaryClientInterceptor {
        return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
                if req == nil {
                        return
                }

		md, _ := metadata.FromIncomingContext(ctx)

		errors := md.Get(ValidationErrorMetaKey)
		if len(errors) > 0 {
			return fmt.Errorf(errors[0])
		}

                return invoker(ctx, method, req, reply, cc, opts...)
        }
}

