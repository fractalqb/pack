package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"git.fractalqb.de/fractalqb/nmconv"
)

//go:generate versioner ../VERSION ./version.go

var mkName = nmconv.Conversion{
	Denorm: nmconv.Camel1Up,
	Norm:   nmconv.Unsep("_"),
}

type att struct {
	key, val string
}

func readInput(filename string) (res []att) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scn := bufio.NewScanner(f)
	for scn.Scan() {
		line := scn.Text()
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}
		sep := strings.IndexRune(line, '=')
		if sep > 0 {
			a := att{key: line[:sep], val: line[sep+1:]}
			res = append(res, a)
		} else {
			log.Fatal("syntax error: ", line)
		}
	}
	return res
}

func writeInput(filename string, attls []att) {
	tmpnm := filename + "~"
	f, err := os.Create(tmpnm)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if f != nil {
			f.Close()
		}
	}()
	for _, a := range attls {
		fmt.Fprintf(f, "%s=%s\n", a.key, a.val)
	}
	err = f.Close()
	f = nil
	if err != nil {
		log.Fatal(err)
	}
	os.Rename(tmpnm, filename)
}

func writeOutput(filename string, atts []att) (bnoChanged []att) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	fmt.Fprintf(f, "package %s\n\nconst (\n", *flagPkg)
	hadBno := false
	for i, a := range atts {
		if a.key == *flagBno {
			n, err := strconv.Atoi(a.val)
			if err != nil {
				log.Fatalf("invalid build number '%s' in attribute %s", a.val, a.key)
			}
			a.val = strconv.Itoa(n + 1)
			atts[i].val = a.val
			hadBno = true
		}
		atnm := mkName.Convert(a.key)
		fmt.Fprintf(f, "\t%s%s = %s\n", *flagPfx, atnm, a.val)
	}
	if !hadBno {
		if len(*flagBno) > 0 {
			a := att{key: *flagBno, val: "1"}
			atts = append(atts, a)
			atnm := mkName.Convert(*flagBno)
			fmt.Fprintf(f, "\t%s%s = %s\n", *flagPfx, atnm, a.val)
		} else {
			atts = nil
		}
	}
	if len(*flagTim) > 0 {
		fmt.Fprintf(f, "\t%s%s = \"%s\"\n", *flagPfx, *flagTim,
			time.Now().Format(time.RFC3339))
	}
	fmt.Fprintln(f, ")")
	return atts
}

func usage() {
	fmt.Fprintf(os.Stderr, "%s v%d.%d.%d%s (%d)\n",
		os.Args[0], Major, Minor, Bugfix, Quality, BuildNo)
	fmt.Fprintln(os.Stderr, "Usage: [flags] input output")
	fmt.Fprintln(os.Stderr, "Flags:")
	flag.PrintDefaults()
}

var (
	flagBno = flag.String("bno", "", "use build number attribute")
	flagPkg = flag.String("pkg", "main", "set the package name for generated file")
	flagPfx = flag.String("p", "", "set name prefix for version constants")
	flagTim = flag.String("t", "", "generate timestamp attribute")
)

func main() {
	flag.Usage = usage
	flag.Parse()
	inNm := flag.Arg(0)
	if len(inNm) == 0 {
		log.Fatal("missing input file argument")
	}
	ouNm := flag.Arg(1)
	if len(ouNm) == 0 {
		log.Fatal("missing output file argument")
	}
	attls := readInput(inNm)
	attls = writeOutput(ouNm, attls)
	if attls != nil {
		writeInput(inNm, attls)
	}
}
