## librespot-go

This is a production-grade adaptation of [librespot-golang](https://github.com/librespot-org/librespot-golang), which itself is an adaptation of [librespot for Rust](https://github.com/librespot-org/librespot) and [librespot-java](https://github.com/librespot-org/librespot-java).

### Objectives
  - Provide core _librespot_ functionality in Go while departing from the constraints of its predecessors.
  - Refactor predecessor into independent components that are Go friendly (e.g. `io.ReadSeekCloser`).
  - Focus on core functionality: peripheral functionality should _consume_ core functionality rather than being _embedded_ within it.

### Points of Interest

  |          |             |
  |----------|-------------|
  | [examples](https://github.com/art-media-platform/librespot-go/tree/main/librespot/examples)                  | PRs welcome               |
  | [api.respot.go](https://github.com/art-media-platform/librespot-go/blob/main/librespot/respot/api.respot.go) | package entry points      |
  | [api.media.go](https://github.com/art-media-platform/amp.SDK/blob/main/stdlib/media/api.media.go)            | media data asset support  |

### Contributing

Contributions aligned with the above objectives are welcome. As this repository is in production, pull requests should demonstrate clear benefit. If you have ideas for improvements, please start a discussion.

