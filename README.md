## go-librespot

[![Build Status](https://github.com/amp-3d/amp-librespot-go/workflows/go/badge.svg?branch=main)](https://github.com/amp-3d/amp-librespot-go/actions)


This Go package is an adaption of [librespot-golang](https://github.com/librespot-org/librespot-golang), which itself is an adaption of a [librespot for Rust](https://github.com/librespot-org/librespot) and [librespot-java](https://github.com/librespot-org/librespot-java).


### Why this fork?
  - Offer essential _librespot_ functionality in Go while departing from the constraints of its predecessor.
  - Refactor its predecessor into components that are independent from each other and  use standard Go idioms and interfaces.
  - Focus on core functionality vs. peripheral functionality.  For multiple reasons,  peripheral components should consume core functionality, not be embedded alongside it.

I will happily support efforts to merge the work done here with [librespot-golang](https://github.com/librespot-org/librespot-golang), but it will require the support of others who are in alignment with the whys above.  If you are interested in this or contributing to this repo, then please start or join a discussion here.
