# pack
[![GoDoc](https://godoc.org/github.com/fractalqb/pack?status.svg)](https://godoc.org/github.com/fractalqb/pack)

`import "git.fractalqb.de/fractalqb/pack"`

---
# Intro
For Go packages it is common and good style to have packages installabel with
a simple `go get...` command. This is supported in a great way by having so many
good tools being part of the Go distribution. However sometimes applications
might require to be distributed as a binary package with a crafted distribution
file-tree – especially when users are not expected to compile the application
themselves.

To have a portable way to pack such a distribution one would need a portable 
runtime for the program that does the packing. So why not write the packing
program in Go? Most of the batteries are already included in Go but IMHO they 
could use a little polishing,  i.e. lifting the API abstraction, to make the
packing code even more simple. And exactly that's what this package is provided
for!

## By the way…
There is just a little tool included to create a source file with constants
defining the current version of the package. The idea is extraordinary simple:
Write down the version parameters in a separate text file from which not only
the Go code can be generated but that likely can be processed by a lot of other
programs too, e.g. shell, makefiles, Python:

```shell
major=0
minor=2
bugfix=0
quality="a"
build_no=13
```

With such a simple file an one line, e.g. in [versioner.go]


# Install
To get the library for writing your packing program:

`go get -u git.fractalqb.de/fractalqb/pack`

To get the `versioner` binary:

`go install -u git.fractalqb.de/fractalqb/pack/versioner`