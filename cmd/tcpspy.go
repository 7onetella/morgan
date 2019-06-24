package cmd

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/spf13/cobra"
)

type redisLogger struct {
	debug bool
}

func (r redisLogger) Println(v ...interface{}) {
	if r.debug {
		log.Println(v...)
	}
}

func (r redisLogger) Printf(msg string, v ...interface{}) {
	if r.debug {
		log.Printf(msg, v...)
	}
}

var tcpSpyCmdDebug bool
var tcpSpyCmdSpyMode bool
var tcpSpyCmdRemoteTimeout int64 = 60
var logger redisLogger

var tcpSpyCmd = &cobra.Command{
	Use:     "tcpspy <address1>",
	Short:   "Spy on tcp connection",
	Long:    `Spy on tcp connection`,
	Example: "replicate remote-server:5000",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println()

		address1 := args[0]

		l, err := net.Listen("tcp", "localhost:5000")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		Println("starting tcp spy")

		for {
			conn, err := l.Accept()
			if err != nil {
				log.Fatal(err)
			}

			go handleRequest(conn, address1)
		}

	},
}

func init() {
	rootCmd.AddCommand(tcpSpyCmd)

	flags := tcpSpyCmd.Flags()

	flags.BoolVarP(&tcpSpyCmdDebug, "debug", "d", false, "debug redis replicate")
	flags.BoolVarP(&tcpSpyCmdSpyMode, "spy", "s", false, "intercept the command and show response from each call")
	flags.Int64Var(&tcpSpyCmdRemoteTimeout, "timeout", 60, "read timeout for redis replicate")
	logger = redisLogger{
		debug: tcpSpyCmdDebug,
	}
}

func handleRequest(conn net.Conn, addr1 string) {
	logger.Println("start of communication")
	defer conn.Close()

	upstream, err := dialTCP(addr1)
	if err != nil {
		log.Println(err)
		return
	}
	defer upstream.Close()

	// sending input to upstream. don't block here, do async
	var w io.Writer
	w = upstream
	go io.Copy(w, conn)

	defer logger.Println("end of communication") // the following blocks

	// Block and send server response downstream
	// send downstream : sending output to downstream
	for {
		upstream.SetReadDeadline(time.Now().Add(time.Duration(tcpSpyCmdRemoteTimeout) * time.Second))
		var buf [128]byte
		n, err := upstream.Read(buf[:])
		if err != nil {
			logger.Printf("upstream %d bytes copied, err = %v", n, err)
			return
		}

		// intercept and do something
		log.Print(string(buf[:n]))

		_, err = conn.Write(buf[:n])
		if err != nil {
			logger.Println(err)
		}
	}
}

// tcp is assumed
func dialTCP(address string) (net.Conn, error) {
	upstream, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	upstream.SetReadDeadline(time.Now().Add(time.Duration(tcpSpyCmdRemoteTimeout) * time.Second))
	return upstream, err
}
