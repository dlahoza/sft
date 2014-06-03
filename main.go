/*

Simple File Transfer

Copyright (c) 2014 Dmitry Lagoza
Use of this source code is governed by a MIT
license that can be found in the LICENSE file.

WWW: http://lagoza.name/
Email: dmitry@lagoza.name

*/

package main

import (
	"fmt"
	flag "github.com/dotcloud/docker/pkg/mflag"
	"io"
	"net"
	"os"
	"strconv"
)

var (
	port               int
	s, d, h, verbose   bool
	filename, hostname string
)

func PrintDefaults() {
	fmt.Fprintf(os.Stderr, "usage: %s [[-s|--source]|[-d|--destination]] [-p num|--port=num] [-h|--help] [HOSTNAME] FILENAME\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "  HOSTNAME\tserver hostname or IP to connect to")
	fmt.Fprintln(os.Stderr, "  FILENAME\tname of file to send or recieve")
}

func init() {
	flag.Usage = PrintDefaults
	flag.BoolVar(&s, []string{"s", "-source"}, false, "start as source server (send FILENAME to client)")
	flag.BoolVar(&d, []string{"d", "-destination"}, false, "start as destination server (recieve FILENAME to client)")
	flag.IntVar(&port, []string{"p", "-port"}, -1, "use port")
	flag.BoolVar(&h, []string{"h", "-help"}, false, "display this help")
	flag.BoolVar(&verbose, []string{"v", "-verbose"}, false, "be verbose")
	flag.Parse()
	if h {
		PrintDefaults()
		os.Exit(0)
	}
	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "Please specify the file!")
		PrintDefaults()
		os.Exit(1)
	}
	if d && s {
		fmt.Fprintln(os.Stderr, "You can use source OR destination flag, not both!")
		PrintDefaults()
		os.Exit(1)
	}
	if (d || s) && flag.NArg() > 1 {
		fmt.Fprintln(os.Stderr, "For server mode don't specify HOSTNAME! You can use -d or -s or HOSTNAME.")
		PrintDefaults()
		os.Exit(1)
	}
	if (!d && !s) && flag.NArg() < 2 {
		fmt.Fprintln(os.Stderr, "For client mode specify HOSTNAME!")
		PrintDefaults()
		os.Exit(1)
	}

}

func Server() {
	filename = flag.Arg(0)
	if port == -1 {
		port = 18000
	}
	if d {
		if verbose {
			fmt.Fprintln(os.Stderr, "Mode: Destination\nPort: "+strconv.Itoa(port)+"\nFile: "+filename)
		}

		ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
		if err != nil {
			os.Exit(1)
		}
		defer ln.Close()
		if verbose {
			fmt.Fprint(os.Stderr, "Waiting for connection... ")
		}
		conn, err := ln.Accept()
		if err != nil {
			os.Exit(1)
		}
		defer conn.Close()
		if verbose {
			fmt.Fprintln(os.Stderr, "connected from "+conn.RemoteAddr().String())
		}
		var f *os.File
		if filename == "-" {
			f = os.Stdout
		} else {
			f, err = os.Create(filename)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error while file creation:", err)
				os.Exit(1)
			}
		}
		defer f.Close()
		fmt.Fprint(conn, "D")
		written, err := io.Copy(f, conn)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error while Copy:", err)
			os.Exit(1)
		}
		if verbose {
			fmt.Fprintln(os.Stderr, "Recieved "+strconv.Itoa(int(written))+" bytes")
		}
	}
	//Source
	if s {
		if verbose {
			fmt.Fprintln(os.Stderr, "Mode: Source\nPort: "+strconv.Itoa(port)+"\nFile: "+filename)
		}

		ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
		if err != nil {
			os.Exit(1)
		}
		defer ln.Close()
		if verbose {
			fmt.Fprint(os.Stderr, "Waiting for connection... ")
		}
		conn, err := ln.Accept()
		if err != nil {
			os.Exit(1)
		}
		defer conn.Close()
		if verbose {
			fmt.Fprintln(os.Stderr, "connected from "+conn.RemoteAddr().String())
		}
		var f *os.File
		if filename == "-" {
			f = os.Stdin
		} else {
			f, err = os.Open(filename)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error while file opening:", err)
				os.Exit(1)
			}
		}
		defer f.Close()
		fmt.Fprint(conn, "S")
		written, err := io.Copy(conn, f)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error while Copy:", err)
			os.Exit(1)
		}
		if verbose {
			fmt.Fprintln(os.Stderr, "Sent "+strconv.Itoa(int(written))+" bytes")
		}
	}
}

func Client() {
	hostname = flag.Arg(0)
	filename = flag.Arg(1)
	if port == -1 {
		port = 18000
	}
	if verbose {
		fmt.Fprintln(os.Stderr, "Client mode")
	}
	fmt.Fprintln(os.Stderr, "Port: "+strconv.Itoa(port)+"\nFile: "+filename)
	fmt.Fprint(os.Stderr, "Connecting... ")
	conn, err := net.Dial("tcp", hostname+":"+strconv.Itoa(port))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error while file connecting:", err)
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Fprintln(os.Stderr, "ok")
	mode := make([]byte, 1)
	conn.Read(mode)
	switch mode[0] {
	default:
		fmt.Fprintln(os.Stderr, "Unknown server mode")
		os.Exit(1)
	case 'D':
		var f *os.File
		if filename == "-" {
			f = os.Stdin
		} else {
			f, err = os.Open(filename)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error while file open:", err)
				os.Exit(1)
			}
		}
		defer f.Close()
		fmt.Fprintln(os.Stderr, "Mode: Destination")
		io.Copy(conn, f)
	case 'S':
		var f *os.File
		if filename == "-" {
			f = os.Stdout
		} else {
			f, err = os.Create(filename)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error while file creation:", err)
				os.Exit(1)
			}
		}
		defer f.Close()
		fmt.Fprintln(os.Stderr, "Mode: Source")
		io.Copy(f, conn)
	}
}

func main() {
	if d || s {
		Server()
	} else {
		Client()
	}
}
