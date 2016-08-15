package main

// Swaths of this stolen from http://blog.scalingo.com/post/105010314493/writing-a-replacement-to-openssh-using-go-22

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"unsafe"

	"github.com/kr/pty"

	"golang.org/x/crypto/ssh"
)

func main() {
	// An SSH server is represented by a ServerConfig, which holds
	// certificate details and handles authentication of ServerConns.
	sshConfig := &ssh.ServerConfig{
		PublicKeyCallback: func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			if conn.User() != "brimstone" {
				return nil, fmt.Errorf("nope")
			}
			log.Println(key)
			return nil, nil
		},
	}

	private, err := ssh.ParsePrivateKey([]byte(`-----BEGIN RSA PRIVATE KEY-----
MIIJKQIBAAKCAgEApuq0B6alVCp1fI2tnoPdRQWZMXE8p2xfrpDRvgVRM1pGeAnd
XDE8XZcE4ESIAP67ti6oWjvvAFFLP5w3OdIH+Q4a9xQXxFCw4apttleST3Wu953x
pdJdp+xmsTKRfFwaQJC7QPgj0wOFQcI25NC/LHIDqAdX+Irk+4iRl+gVcD+5vmTw
TDo7djzfVfUZkCShfw4k2E8sEcwT5Ya4B+KuY/OodNyOmBgt7Ep3Wp7kqTHVhhum
1ENGhPXR6AXQcYcplTldfk2k9G4sRWWvvNpFqH2dVfSwEyQT2pLCfrZKA239nn6C
wn/7JPXkA3zdU3o/txY0B/jpoxSV7GZz09T+xLY/nqr8Cg78rp5Oy69Znp5WTlEI
y7lf7tpt0k/oNe7LN39De2H+oirfpdQaiyKb3X8gOst4wpLsa/h9kVptRGzFGSW6
bXmRaosMQTAv5bvhvrLMRHuS6c7IA4zdVJwzJDRad+dx0MTsjrK6ShtcQzKTCBkK
wP6Ls7Y6gpRLtXJSjR0bZjPu5foBCqSzBn5cGZxitcFrBO/x1PFgiMUQ7smG2+xb
KxMo1etNKruO5B1j8UNEBFrgcaD67q2P0aAY9512GIt0ByhWSaN4M0w6lqxTyAQU
Kpk7BDSzMgl3tMzWMSBo7B6nkENZ+kUV6tCGqLWh1ZlMgTHyNLfJ39QosmMCAwEA
AQKCAgAIYxCqtb/m6789g+zuYxfSKQbaiiMPEo34OoSfdKrw1p9l1rENudebqEPx
dOAUlLgf3lZNOme2717Fknbf6+LEq+XE9nh/P8KzhBnBKMNMRNCG1qPWvixAjMtY
Kf9PbV1QUzVlfVJnfuzKMhUKCEci15PBdKUB8xCwZttR87JoEnulynKckex883AR
ZKBlMsH+nVpSmB/RwRxa0xsaIlS65vpW7OIpWEWucstufT9mFP/ynh8S2VKIycVD
UE959N2tBtXgy8v3EDYfQo3DAoCvh6hJMmNKguyQdgFZ1pT/eR5eQWMWnNGkuMo9
Np/0Wtcqvu3cXLB4pkcmOrulMDWg3e8L4txMAl/ywY4648CfkSbRDBz5Xeyu+uI6
cWwwoEDoXQvIpO2no6dHDx/C+dtJHH1k8g0wL2a2CtwbMSSxPXjxp6kWVncCwQ9M
BoRh9lZ88EQ9e/sZfp3cDMD9PNtYYhDgLQ4CgUU43A+cd4057codvkxG4OiGGXOy
UKvd33qbBYo89mE8i4vdI7Hrl2hPSsajpghpDy7ZiDnO493kMtPDxxwD6cmpTB88
V93jNC/H8j8MCcc24T6aUekfyGpU00rVQHgqkt5YIRE298iwx3zEOow9N1VXamYq
k6cTWp3aNUviA58GM6QdZczLXaznkmaDQHT3mjD9O4qLLujCAQKCAQEA2RKU3QJ5
W3p6A5fE7hvs7J3iRZNpxY22Vh2LWgDGKEbd6FHbKgq6n8oxQLXX/f+i/P5YX71n
QH4iZlUt9YYpYtK4jg94xEzVPsvoOIITuE4mGWtxI9bvtmEA2VUeODq97zzNptY6
wHj79RbZab8Ntj+XXhwMYJoY3Trh8hqu/ZhjI2t+s78lIYpFi2s1yG1SWFBRJx7l
q+Y0QH9ilj2EJcYZJnj5KTaoBYfpd/VCufbdqtoTkguE94bmzl5catBXvt3/ciup
xjHmb2KgK4tvyhfoLv4l78dCvzD2LCK8HO1XtTtgu87tjZxC+k6MBfA+MMXhT/aq
7QBoyPVDrtGjYwKCAQEAxNmPDy9fNLeGgPveuELM19x6bqdpCOKNwRcwrRbusQTi
koBN1fnsko2qL4KzLDVhUYFQ7vKcp/QfUDjdgV1lqDZj908RCl3o3RQHWrBZOgmp
pLbnYfNeewcawNKcZfI9RMjorvwslEr7g1/DyTGAMXVLOh1WEebMXs5bpiRjj/ob
Ca53Rd2WVprxmcf8w59ts7Z4c7vAJJFJOhC3qTTl9U5HkFYawy++R18RSkSTn12E
oM9yn/+NjKbUSjUirvBewIj5W2hQFkYMpWdHBQj575AW7RJOFRRjLgUSwPvNaNar
7cZvq9d7tctRTZABizU4lfSgBJh0xSuuMdw//OtlAQKCAQEAvvnEEjUtA7hbRHHw
BfR3myzEd6XbtryYoGbZxNNNgv3mGZB6myBZVF+UuXPClWqkwKQcqA6AmpLePN3P
02S3YIQ3bnRaMYnPSrImmiCGrO3EQzTtkzR0LSZmks56CcpUc/gwjgmIIvHN6bCy
koBN1ftYdqmCkjLAkVJOmquzLNU202CtVgJi7oEx8hjednkHqz7uRta+BWBAtEXe
PIPekUUZt2lS/FljtYn/c1RJ4kY4eynlceXEk+kRgpouAKNPr1KsfIvBj+cy7uf8
NpM0RL6HvWBHjA7owECZM/dTPLfrJD0bO+AvyxihLIqSUl8st8hAFBPWCTEE/1gY
teMmNQKCAQEAuAJe1mUD7DV6R+wpv7jB3y68S9+2MZYFyG+zEroTGeplGWlSWSks
2boPUiYs8rBbbmMhkpu7kMyE0Oq3NIxn0Jw3SiDg3v32BWMJlN6wKa6Ko+xN0qQ6
t2pmucSmai3M1BWyXJBh46VMAvxr+hCJsrHgRkzR/h3vANiJl38Air/SsnQiDm8a
b46bNZNaVksbsxho9FaXQBeHif3CkStforUv9F5o2fgxOGpHsVL1Y815gxEoJyQt
30K7wzp0V29eQ0BHSAj0hD2q6JroKm6/pA7fP3ETCGVsaMJZS4iV9OBnCvepv8rI
W6HZOFXa+5Qedx6aznDtBrrUNMucDQLGAQKCAQBpDqM4xaQYA8Te0pcnczJC/+lR
ubWbOE6pH0l34rKHbZIRb8tqhI6PcI6WoozsPNfAPlAaTyIh/NGVEhhhYY7+hLi4
D4VU5fq9XoH+PeZ3gduR05EvfI9SZ2QsoMmEIYLIDFckCkAQCEVILSFFlTQ4SGgI
LMoijoI6Wb1gK11L4WpsuDQabr4a5HDNlyeKwXKMaSD++zv2D16qGx+0uHvLojbi
O/3iwuzeQ8HjgoekdxCVUp0HDn+xnWUQwXfIT4CKfIy+0s97sHQshrzbPV1jL6rY
OL0+DH5Q4FqwmXwG0J4sKluPcXfu+Fq0P8i6YFXa+w3uKjbgczGlCnPjQBcx
-----END RSA PRIVATE KEY-----
`))
	if err != nil {
		panic("Failed to parse private key")
	}

	sshConfig.AddHostKey(private)

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

func handleRequests(reqs <-chan *ssh.Request) {
	for req := range reqs {
		log.Printf("received out-of-band request: %+v", req)
	}
}

// Start assigns a pseudo-terminal tty os.File to c.Stdin, c.Stdout,
// and c.Stderr, calls c.Start, and returns the File of the tty's
// corresponding pty.
func PtyRun(c *exec.Cmd, tty *os.File) (err error) {
	defer tty.Close()
	c.Stdout = tty
	c.Stdin = tty
	c.Stderr = tty
	c.SysProcAttr = &syscall.SysProcAttr{
		Setctty: true,
		Setsid:  true,
	}
	return c.Start()
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

// =======================

// parseDims extracts two uint32s from the provided buffer.
func parseDims(b []byte) (uint32, uint32) {
	w := binary.BigEndian.Uint32(b)
	h := binary.BigEndian.Uint32(b[4:])
	return w, h
}

// Winsize stores the Height and Width of a terminal.
type Winsize struct {
	Height uint16
	Width  uint16
	x      uint16 // unused
	y      uint16 // unused
}

// SetWinsize sets the size of the given pty.
func SetWinsize(fd uintptr, w, h uint32) {
	log.Printf("window resize %dx%d", w, h)
	ws := &Winsize{Width: uint16(w), Height: uint16(h)}
	syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TIOCSWINSZ), uintptr(unsafe.Pointer(ws)))
}

// ByteTo does things with bytes
func ByteTo(data interface{}, byteArray []byte) error {
	log.Println(byteArray)
	buf := bytes.NewReader(byteArray)
	err := binary.Read(buf, binary.BigEndian, data)
	if err != nil {
		return err
	}
	return nil
}
