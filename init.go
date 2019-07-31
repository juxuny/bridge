package bridge

import "flag"

func init() {
	flag.BoolVar(&verbose, "v", true, "display debug output")
	//flag.Parse()
}
