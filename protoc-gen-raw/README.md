# protoc-gen-raw

This generator is used to output raw protobuf data.

## Usage

```bash
protoc *.proto --raw_out=.
```

Running this command will give you a `raw.protobuf` file. You can `proto.Unmarshal` it for debug use.
