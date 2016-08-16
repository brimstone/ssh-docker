package main

// Swaths of this stolen from http://blog.scalingo.com/post/105010314493/writing-a-replacement-to-openssh-using-go-22

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"sync"

	"github.com/kr/pty"

	"golang.org/x/crypto/ssh"
)

func main() {
	// An SSH server is represented by a ServerConfig, which holds
	// certificate details and handles authentication of ServerConns.
	sshConfig := sshConfig()
	// Once a ServerConfig has been configured, connections can be accepted.
	listener, err := net.Listen("tcp4", ":2022")
	if err != nil {
		log.Fatalf("failed to listen on *:2022")
	}

	// Accept all connections
	log.Printf("listening on %s", ":2022")
	for {
		tcpConn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept incoming connection (%s)", err)
			continue
		}
		// Before use, a handshake must be performed on the incoming net.Conn.
		sshConn, chans, reqs, err := ssh.NewServerConn(tcpConn, sshConfig)
		if err != nil {
			log.Printf("failed to handshake (%s)", err)
			continue
		}

		// Check remote address
		log.Printf("new ssh connection from %s (%s)", sshConn.RemoteAddr(), sshConn.ClientVersion())

		// Print incoming out-of-band Requests
		go handleRequests(reqs)
		// Accept all channels
		go handleChannels(chans)
	}
}

func handleChannels(chans <-chan ssh.NewChannel) {
	// Service the incoming Channel channel.
	for newChannel := range chans {
		// Channels have a type, depending on the application level
		// protocol intended. In the case of a shell, the type is
		// "session" and ServerShell may be used to present a simple
		// terminal interface.
		if t := newChannel.ChannelType(); t != "session" {
			newChannel.Reject(ssh.UnknownChannelType, fmt.Sprintf("unknown channel type: %s", t))
			continue
		}
		channel, requests, err := newChannel.Accept()
		if err != nil {
			log.Printf("could not accept channel (%s)", err)
			continue
		}

		// allocate a terminal for this channel
		log.Print("creating pty...")
		// Create new pty
		f, tty, err := pty.Open()
		if err != nil {
			log.Printf("could not start pty (%s)", err)
			continue
		}

		var shell string
		shell = os.Getenv("SHELL")
		if shell == "" {
			shell = "/bin/bash"
		}

		// Sessions have out-of-band requests such as "shell", "pty-req" and "env"
		go func(in <-chan *ssh.Request) {
			for req := range in {
				ok := false
				log.Println("Request:", req.Type)
				switch req.Type {
				case "exec":
					ok = true
					command := string(req.Payload[4 : req.Payload[3]+4])
					log.Println("Command:", command)
					cmd := exec.Command(shell, []string{"-c", command}...)

					cmd.Stdout = channel
					cmd.Stderr = channel
					cmd.Stdin = channel

					err := cmd.Start()
					if err != nil {
						log.Printf("could not start command (%s)", err)
						continue
					}

					// teardown session
					go func() {
						_, err := cmd.Process.Wait()
						if err != nil {
							log.Printf("failed to exit bash (%s)", err)
						}
						channel.Close()
						log.Printf("session closed")
					}()
				case "shell":
					cmd := exec.Command(shell)
					cmd.Env = []string{"TERM=xterm"}
					err := PtyRun(cmd, tty)
					if err != nil {
						log.Printf("%s", err)
					}

					// Teardown session
					var once sync.Once
					close := func() {
						channel.Close()
						log.Printf("session closed")
					}

					// Pipe session to bash and visa-versa
					go func() {
						io.Copy(channel, f)
						once.Do(close)
					}()

					go func() {
						io.Copy(f, channel)
						once.Do(close)
					}()

					// We don't accept any commands (Payload),
					// only the default shell.
					if len(req.Payload) == 0 {
						ok = true
					}
				case "pty-req":
					// Responding 'ok' here will let the client
					// know we have a pty ready for input
					ok = true
					// Parse body...
					// FIXME Dangerous to just pull the 3rd bit
					termLen := req.Payload[3]
					termEnv := string(req.Payload[4 : termLen+4])
					w, h := parseDims(req.Payload[termLen+4:])
					SetWinsize(f.Fd(), w, h)
					log.Printf("pty-req '%s'", termEnv)
				case "window-change":
					w, h := parseDims(req.Payload)
					SetWinsize(f.Fd(), w, h)
					continue //no response
				case "env":
					log.Printf("%#v\n", req)
					payload := req.Payload
					log.Printf("%v %s\n", payload, payload)
					// Get key size
					var keyLen int32
					ByteTo(&keyLen, payload[0:4])
					payload = payload[4:]
					key := string(payload[0:keyLen])
					payload = payload[keyLen:]
					// Get value size
					var valueLen int32
					ByteTo(&valueLen, payload[0:4])
					payload = payload[4:]
					value := string(payload[0:valueLen])
					log.Println(key, ":", value)
					// FIXME, arg, not actually opening a channel to anywhere
					//channel.Setenv(key, value)
					ok = true
				}

				if !ok {
					log.Printf("declining %s request...", req.Type)
				}

				req.Reply(ok, nil)
			}
		}(requests)
	}
}
