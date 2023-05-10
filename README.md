## go-librespot


This Go package is an adaption of [librespot-golang](https://github.com/librespot-org/librespot-golang), which itself is an adaption of a [librespot for Rust](https://github.com/librespot-org/librespot) and [librespot-java](https://github.com/librespot-org/librespot-java).  


Why this fork?
  - Offer essential _librespot_ functionality for Go while departing from the constraints of its predecessor.
  - Refactor its predecessor into proper interfaces that leverage Go.
  - Focus on core functionality vs. peripheral functionality (e.g. audio conversion, remote control).  For multiple reasons,  non-core functionality should be in a consuming repo, not the core repo.

I will happily support efforts to merge the work done here with [librespot-golang](https://github.com/librespot-org/librespot-golang), but it will require the support of others who are in alignment with the whys above.  If you are interested in this or contributing to this repo, then please start or join a discussion here and let's get to work.  