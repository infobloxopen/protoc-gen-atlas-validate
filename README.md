# protoc-gen-atlas-validate

The main purpose of this plugin is to generate a code under pb.atlas.validate.go
that serves several purposes:

  - Ability to configure 'allow unknown fields' on several levels: per method, per service, per proto-file

  - Validate basic types.

  - (TODO) Ability to configure 'read-only' fields

  - (TBD-TODO) Possibly this can be transformed to a full-fledged ad-hoc JSON marshaller
    with per-service/method/file options similar to ones that OpenAPI provides.

## Usage

Include following lines in your `.proto` file:

### Import

```
import "github.com/infobloxopen/protoc-gen-atlas-validate/options/atlas_validate.proto";
```

### Specifying options

Service option:

```
service Groups {
        option (atlas_validate.service).allow_unknown_fields = true;
        rpc Create(Group) returns (EmptyResponse) {
...
```

Method option:

```
        rpc Update(UpdateProfileRequest) returns (EmptyResponse) {
                option (atlas_validate.method).allow_unknown_fields = true;
                option (google.api.http) = {
                        put: "/profiles/{payload.id}";
                        body: "payload";
                };
        }
}
```

Global option:

```
option (atlas_validate.file).allow_unknown_fields = false;
```

### Generation

Note that this plugin heavily relies on patterns generated by protoc-gen-grpc-gateway plugin:

```
protoc -I/usr/local/include \
	-I. -I$GOPATH/src/ \
	-I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway \
	-I./vendor \
	-I$GOPATH/src/github.com/googleapis \
		--grpc-gateway_out="logtostderr=true:$GOPATH/src" \
		--atlas-validate_out="$GOPATH/src" \
			<path-to-your-file>
```

The following will generate pb.atlas.validate.go file that contains validation
logic and MetadataAnnotator that you will have to include in GRPC Server options.

### Usage

Import atlas-validate Interceptor:

```
import atlas_validate "github.com/infobloxopen/protoc-gen-atlas-validate/interceptor"
```

Add generated AtlasValidateAnnotator (from \*.pb.atlas.validate.go) to a Metadata anotators:

```
gateway.WithGatewayOptions(
	runtime.WithMetadata(pb.AtlasValidateAnnotator),
)
```

Add interceptor that extracts error from metadata and returns it to a user:

```
gateway.WithDialOptions(
	[]grpc.DialOption{grpc.WithInsecure(), grpc.WithUnaryInterceptor(
		grpc_middleware.ChainUnaryClient(
			[]grpc.UnaryClientInterceptor{
				gateway.ClientUnaryInterceptor,
				atlas_validate.ValidationClientInterceptor(),
			},
		)
	)
)
```