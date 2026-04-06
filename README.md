# librespot-go

A Go package for [Spotify](https://www.spotify.com/) audio streaming, forked from [librespot-golang](https://github.com/librespot-org/librespot-golang) (itself derived from [librespot](https://github.com/librespot-org/librespot) for Rust and [librespot-java](https://github.com/librespot-org/librespot-java)).

This fork is maintained as a production dependency of [art.media.platform](https://github.com/art-media-platform/amp.SDK) — specifically the `app.av.spotify` module in [amp.planet](https://github.com/art-media-platform/amp.planet).

## Goals

- **Go-native API** — package consumption, not CLI.  Media assets implement `io.ReadSeekCloser` via the amp.SDK [`data.Asset`](https://github.com/art-media-platform/amp.SDK/blob/main/stdlib/data/asset.go) interface.
- **Independent components** — core, mercury, metadata, asset, and spirc are separate packages with clean boundaries.
- **Minimal surface** — core streaming functionality only.  Peripheral features should consume the core, not live inside it.

## Packages

| Package | Purpose |
|---------|---------|
| [`respot/`](librespot/respot/) | Top-level API — session management, authentication, player |
| [`core/`](librespot/core/) | Spotify session, handshake, transport |
| [`asset/`](librespot/asset/) | Audio asset fetching, chunked streaming, decryption |
| [`mercury/`](librespot/mercury/) | Mercury protocol (Spotify's internal pub/sub RPC) |
| [`metadata/`](librespot/metadata/) | Track, album, artist, playlist metadata resolution |
| [`spirc/`](librespot/spirc/) | Spotify Connect remote control protocol |
| [`discovery/`](librespot/discovery/) | mDNS/SSDP device discovery |
| [`Spotify/`](Spotify/) | Generated protobuf types (protoc-gen-go v1) |

## Quick Start

```go
import "github.com/art-media-platform/librespot-go/librespot/respot"
```

See [`librespot/examples/`](librespot/examples/) for usage patterns.

## Note on Protobufs

The `Spotify/` directory contains protobuf types generated with `protoc-gen-go` v1 (`github.com/golang/protobuf`).  Code that marshals/unmarshals these types must use the v1 proto package, not `google.golang.org/protobuf/proto`.

## License

[MIT](LICENSE) — original copyright Paul Lietar, badfortrains, Guillaume Lesniak.
