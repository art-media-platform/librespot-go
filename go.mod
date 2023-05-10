module github.com/arcspace/go-librespot

go 1.18

// replace github.com/arcspace/go-cedar => ../go-cedar
replace github.com/arcspace/go-arcspace => ../go-arcspace

require (
	github.com/arcspace/go-arcspace v0.0.0-20230426065400-4ccbafbaca4e
	github.com/arcspace/go-cedar v1.2023.1
	github.com/golang/protobuf v1.5.3
	golang.org/x/crypto v0.8.0
	google.golang.org/protobuf v1.30.0
)

require (
	github.com/brynbellomy/klog v0.0.0-20200414031930-87fbf2e555ae // indirect
	github.com/h2non/filetype v1.1.3 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rs/cors v1.9.0 // indirect
	golang.org/x/sync v0.2.0 // indirect
)
