# lib-post-interchange

**lib-post-interchange** is a repo to house data format libraries for video post-production interchange formats.

This includes:
* ALE (Avid Log Exchange) `.ale` -> libale
* CDL (Color Decision List) `.cdl` (et al) -> libcdl
* EDL (Edit Decision List) `.edl` -> libedl

It is designed to support any programming language but currently supports:

* Go (future)
* JavaScript
* Python

Also planned is a CLI conversion utility. It will function as a convenient CLI tool to take input files in ALE, CDL or EDL and output files in ALE, CDL or EDL as well as other data formats including CSV.

It will also serve as a reference implementation for the libraries in Go.



## Documentation

### JavaScript

See [javascript/README.md](javascript/README.md).

### Python

See [python/README.md](python/README.md) and [python/examples](python/examples).

## Development

February 2025: Under development.

### Naming

Perhaps this project could have a shorter name?

What to do about libxxx where that is already used? I.e. if `libale` is already used, `libavidle`? `libasccdl`, or `libeditdl`.

### TODO

* Upload it to pip! Installation and usage should be as simple as:
```bash
pip install libale
```

* Develop the CLI conversion utility

* Gather a whole heap of real world ALEs to iron out any major things that I'm missing

* Implement unit testing

