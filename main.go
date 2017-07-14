package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"github.com/moxtsuan/go-nkf/nkf"
)

var (
	inJIS     = flag.Bool("J", false, "Input JIS(ISO2022JP)")
	inEUC     = flag.Bool("E", false, "Input EUCJP")
	inSJIS    = flag.Bool("S", false, "Input Shift_JIS")
	inUTF8    = flag.Bool("W", false, "Input UTF-8")
	outJIS    = flag.Bool("j", false, "Output JIS(ISO2022JP)")
	outEUC    = flag.Bool("e", false, "Output EUCJP")
	outSJIS   = flag.Bool("s", false, "Output Shift_JIS")
	outUTF8   = flag.Bool("w", false, "Output UTF-8")
	unix      = flag.Bool("Lu", false, "LF(UNIX) Newline")
	windows   = flag.Bool("Lw", false, "CRLF(Windows) Newline")
	macintosh = flag.Bool("Lm", false, "CR(Macintosh) Newline")
	override  = flag.Bool("override", false, "Override")
	guess     = flag.Bool("g", false, "Detect Charset")
)

func setIn() string {
	switch {
	case *inJIS:
		return "ISO2022JP"
	case *inEUC:
		return "EUCJP"
	case *inSJIS:
		return "ShiftJIS"
	case *inUTF8:
		return "UTF8"
	default:
		return ""
	}
}

func setOut() string {
	switch {
	case *outJIS:
		return "ISO2022JP"
	case *outEUC:
		return "EUCJP"
	case *outSJIS:
		return "ShiftJIS"
	case *outUTF8:
		return "UTF8"
	default:
		return "ISO2022JP"
	}
}

func setNl() string {
	switch {
	case *unix:
		return "UNIX"
	case *windows:
		return "WINDOWS"
	case *macintosh:
		return "MACINTOSH"
	default:
		return ""
	}
}

func main() {
	var fp *os.File
	var err error
	flag.Parse()
	if len(flag.Args()) < 1 {
		fp = os.Stdin
	} else {
		fp, err = os.Open(flag.Args()[0])
		if err != nil {
			log.Fatal(err)
		}
		defer fp.Close()
	}

	in := setIn()
	out := setOut()
	nl := setNl()
	if *guess {
		charset, _ := nkf.Guess(fp)
		fmt.Println(charset)
	} else {
		str, err := nkf.Convert(fp, in, out, nl)
		if err != nil {
			log.Fatal(err)
		}
		if *override && !(len(flag.Args()) < 1) {
			fp.Close()
			fp, err = os.OpenFile(flag.Args()[0], os.O_WRONLY, 0666)
			if err != nil {
				log.Fatal(err)
			}
			defer fp.Close()
			fmt.Fprint(fp, str)
		} else {
			fmt.Fprint(os.Stdout, str)
		}
	}
}
