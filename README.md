## librespot-go

This Go package is an adaption of [librespot-golang](https://github.com/librespot-org/librespot-golang), which itself is an adaption of a [librespot for Rust](https://github.com/librespot-org/librespot) and [librespot-java](https://github.com/librespot-org/librespot-java).

### Why this fork?
  - Replace _librespot_ functionality in Go while departing from constraints of its predecessor.
  - Refactor predecessor into components that are independent from each other and use standard Go idioms.
  - Focus on core functionality vs. peripheral functionality: peripheral functionality should _consume_ core functionality, not be _embedded_ within it.

I will happily support efforts to merge the work done here with [librespot-golang](https://github.com/librespot-org/librespot-golang), but it requires the support of others who are in alignment with the whys above.  If you are interested in this or contributing to this repo, then please start or join a discussion here.