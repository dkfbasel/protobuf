# Protobuf helpers for Golang

This directory `types` contains custom proto messages and corresponding go types
that make working with null values in a database easier.

In addition, the utility `protoc-go-tags` can be used to specify additional
go struct tags in the proto definitions. This can be used to specify database
column names or validation options.

The utility will not yet replace existing json tags on the proto structs. Please
note that we actually parse the go code and modify the AST (abstract syntax tree)
to add your additional go tags. This should be somewhat more reliable than
regular expressions alone.

```
// proto definitions
message StarfleetShip {
	string name = 1;

	// use a different db column name for the departure time
	// `db:"we_are_leaving_at"`
	dkfbasel.protobuf.Timestamp departure_time = 4;
}

// go struct definition
type StarfleetShip struct {
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty" `

	// use a different db column name for the departure time
	// `db:"we_are_leaving_at"`
	DepartureTime *dkfbasel_protobuf.Timestamp `protobuf:"bytes,4,opt,name=departure_time,json=departureTime" json:"departure_time,omitempty" db:"we_are_leaving_at"`
}
```

To add the additional tags to your structs proceed as following - a complete
example can be found in the directory `example`.

1. Add the tag as comment above the respective property. Unfortunately the tag
can currently not be placed on the same line, as the protoc grpc plugin will
strip out all comments that are on the same line.

2. Run the protoc generator with the grpc plugin.

3. Run the protoc-go-tags utility on the generated files. You need to specify
the directory that the utility should process and it will recursively scan
all files and transform all go files with the ending .pb.go

4. The additional tags should now be included in your .pb.go files

# SQL-Queries with Named Parameters

We are often using sqlx to simplify database access for our programs and found
that the matching of named parameters to struct properties would use lowercase
name transforms and thus not work with our snakecased column names.

It is however trivial to change the name mapping by simply passing a custom
mapper function to the slqx db connection. An example of this can be found
in `example/main.go`.
