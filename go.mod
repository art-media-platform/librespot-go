module github.com/art-media-platform/amp-librespot-go

go 1.22

// replace github.com/art-media-platform/amp.SDK => ../amp.SDK

require (
	github.com/art-media-platform/amp.SDK v0.8.7
	github.com/golang/protobuf v1.5.4
	golang.org/x/crypto v0.31.0
	google.golang.org/protobuf v1.36.1
)

require (
	github.com/brynbellomy/klog v0.0.0-20200414031930-87fbf2e555ae // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect; indirectgo mod tidy
)
