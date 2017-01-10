package main

import (
	"os"
	"fmt"
	"time"
	"flag"
	"encoding/json"
	"github.com/miekg/dns"
)

var DEBUG = false
var VERBOSE = false
var TIMEOUT_DEFAULT = 2000
var TIMEOUT_MIN = 10
var TIMEOUT_MAX = 30000

func usage() {
	fmt.Printf("Usage: dnsquery [server] [target]\n")
	os.Exit(1)
}

func errprint(msg string) {
	print(msg)
	os.Exit(1)
}

func print(msg string) {
	t := time.Now().Local()
	output := fmt.Sprintf("%s %s", t.Format("20060102_15:04:05"), msg)
	fmt.Print(output)
}

/*
func arguments(args []string) (string, string) {
	numArgs := len(args)
	if numArgs < 3 {
		usage()
		os.Exit(1)
	}

	return args[1], args[2]
}
*/

/*
 *
 *	MAIN
 *
 */

func main() {
	//var server, target = arguments(os.Args)

	var debugArg   = flag.Bool("debug", false, "enable debug-output")
	var verboseArg = flag.Bool("verbose", false, "enable verbose output")
	var serverArg  = flag.String("server", "none", "nameserver to query")
	var targetArg  = flag.String("target", "none", "name to query")
	var typeArg    = flag.String("type", "a", "recordtype to query for")
	var timeoutArg = flag.Uint("timeout", 2000, "timeout in ms")
	flag.Parse()
	if *debugArg == true {
		print("DEBUG: activating debugmode..\n")
		DEBUG = true
	}
	if *verboseArg == true {
		VERBOSE = true
	}
	var server = *serverArg
	if server == "none" {
		errprint("ERROR: no server given.\n")
		os.Exit(1)
	}
	var target = *targetArg
	if target == "none" {
		errprint("ERROR: no target given.\n")
		os.Exit(1)
	}
	var rrtype uint16
	switch *typeArg {
		case "a":
			rrtype = dns.TypeA
		case "any":
			rrtype = dns.TypeANY
		case "mx":
			rrtype = dns.TypeMX
		case "ns":
			rrtype = dns.TypeNS
		default:
			errprint(fmt.Sprintf("ERROR: unknown type: %s\n", *typeArg))
			os.Exit(1)
	}
	var timeout = int(*timeoutArg)
	if timeout < TIMEOUT_MIN || timeout > TIMEOUT_MAX {
		errprint(fmt.Sprintf("ERROR: invalid timeout-argument (%d), using the default (%d).\n", timeout, TIMEOUT_DEFAULT))
	}

	c := dns.Client{}
	c.Timeout = time.Duration(timeout) * time.Millisecond
	m := dns.Msg{}
	m.SetQuestion(target + ".", rrtype)
	r, t, err := c.Exchange(&m, server + ":53")
	if err != nil {
		errprint(fmt.Sprintf("Exchange(): %s@%s %.2f ms ERROR:%s\n", target, server, t.Seconds() * 1000.0, err.Error()))
	}

	if DEBUG == true {
		out, _ := json.Marshal(r)
		fmt.Println(string(out))
	}

	var rcode string
	switch r.MsgHdr.Rcode {
		case 0:
			//rcode = "NOERROR"
			rcode = "OK"
		case 1:
			rcode = "FORMERR"
		case 2:
			rcode = "SERVFAIL"
		case 3:
			rcode = "NXDOMAIN"
		case 4:
			rcode = "NOTIMP"
		case 5:
			rcode = "REFUSED"
		case 6:
			rcode = "YXDOMAIN"
		case 7:
			rcode = "XRRSET"
		case 8:
			rcode = "NOTAUTH"
		case 9:
			rcode = "NOTZONE"
		default:
			rcode = "UNKNOWN RCODE!"
	}
	//print(fmt.Sprintf("%s\n", rcode))

	/***
	var qstatus string
	if len(r.Answer) == 0 {
		qstatus := "ERROR: EMPTY_RESPONSE"
		os.Exit(1)
	}
	***/

	/***
	qstatus := "OK"
	print(fmt.Sprintf("RESULT: %s@%s %.2f ms %s\n", target, server, t.Seconds() * 1000.0, qstatus))
	***/

	if r.MsgHdr.Rcode != 0 {
		print(fmt.Sprintf("RESULT: %s@%s %.2f ms %s\n", target, server, t.Seconds() * 1000.0, rcode))
		os.Exit(1)
	}
	
	print(fmt.Sprintf("RESULT: %s@%s %.2f ms %s\n", target, server, t.Seconds() * 1000.0, rcode))

	if VERBOSE {
		print("[ANSWER]\n")
		for _, rr := range r.Answer {
			print(fmt.Sprintf("%s\n", rr.String()))

			/*
			var header = rr.Header()
			var name = header.Name
			var rrtype = header.Rrtype
			var ttl = header.Ttl
			print(fmt.Sprintf("name=%s rrtype=%d ttl=%d\n", name, rrtype, ttl))
			*/
		}

		print("[EXTRA]\n")
		for _, rr := range r.Extra {
			print(fmt.Sprintf("%s\n", rr.String()))
		}
	}

	/***
	// ans = RR (dns.go)

	// type A struct {
    //   Hdr RR_Header
    //   A net.IP `dns:"a"`
	// }

	for _, ans := range r.Answer {
		//Arecord := ans.(*dns.A)
		//log.Printf("%s", Arecord.A)
		errprint(fmt.Sprintf("%s\n", ans.String())
	}
	***/
}
