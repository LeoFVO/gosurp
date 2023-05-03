package smtp

import (
	"net"
	"strings"

	log "github.com/sirupsen/logrus"
)


func (server Server) Listen() error {
	log.Infof("Server started listening on %s:%s", server.Hostname, server.Port)
	// The ln variable represents the listener object, which is used to accept incoming connections. 
	ln, err := net.Listen("tcp", server.Hostname + ":" + server.Port)
	if err != nil {
		return err
	}
	// The defer ln.Close() statement ensures that the listener is closed when the function exits.
	defer ln.Close()

	// Accept incoming connections
	// When a new connection is accepted, the handleConnection() function is called to handle incoming messages:
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Errorln(err)
			continue
		}

		// Handle incoming connection
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	// Send initial greeting message according to the RFC
	conn.Write([]byte("220 localhost ESMTP Service Ready\r\n"))
	log.Debugf("Connection from %s", conn.RemoteAddr().String())

	// Buffer for incoming messages
	buf := make([]byte, 1024)

	// Loop to handle incoming commands
	for {
		// Read incoming message
		n, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
			break
		}

		// Parse incoming command
		cmd := string(buf[:n]) // Convert bytes to string
		parts := strings.Split(cmd, " ") // Split string into parts
		command := strings.TrimSpace(parts[0]) // Extract command
		args := parts[1:] // Extract arguments

		// Handle incoming command
		switch command {
			case "EHLO", "HELO":
				if len(args) < 1 {
					conn.Write([]byte("501 Syntax: HELO hostname\r\n"))
					continue
				}
				log.Tracef("Received command %s %+v", command, string(args[0]))
				/*
				* The list of features sent in the code includes:
				*
				* SIZE 31457280: This indicates that the server supports messages up to 31457280 bytes in size.
				* 8BITMIME: This indicates that the server supports 8-bit MIME encoding.
				* STARTTLS: This indicates that the server supports the STARTTLS command, which can be used to initiate a secure TLS connection.
				* The server may support other features that can be included in the list. The purpose of sending the list of features is to allow the client to know what features are available and to negotiate which features to use.
				*/
				conn.Write([]byte("250 localhost\r\n"))
				conn.Write([]byte("250-SIZE 31457280\r\n"))
				conn.Write([]byte("250-8BITMIME\r\n"))
				conn.Write([]byte("250-STARTTLS\r\n"))
				conn.Write([]byte("250-OK\r\n"))
			case "MAIL":
				if len(args) < 1 {
					log.Tracef("Received command %s %+v", command, args)
					conn.Write([]byte("501 Syntax: MAIL FROM:<address>\r\n"))
					break
				}
				if !strings.HasPrefix(string(args[0]), "FROM:") {
					log.Tracef("Received command %s %+v", command, args)
					conn.Write([]byte("501 Syntax: MAIL FROM:<address>\r\n"))
					break
				}
				log.Tracef("Received command %s %+v", command, string(args[0]))
				conn.Write([]byte("250 OK\r\n"))
			case "RCPT":
				if len(args) < 1 {
					log.Tracef("Received command %s %+v", command, args)
					conn.Write([]byte("501 Syntax: RCPT TO:<address>\r\n"))
					break
				}
				if !strings.HasPrefix(string(args[0]), "TO:") {
					log.Tracef("Received command %s %+v", command, args)
					conn.Write([]byte("501 Syntax: RCPT TO:<address>\r\n"))
					break
				}
				log.Tracef("Received command %s %+v", command, string(string(args[0])))
				conn.Write([]byte("250 OK\r\n"))
			case "DATA":
				if len(args) > 0 {
					log.Tracef("Received command %s %+v", command, args)
					conn.Write([]byte("501 Syntax: DATA\r\n"))
					break	
				}
				log.Tracef("Received command %s", command)
				conn.Write([]byte("354 Start mail input; end with <CRLF>.<CRLF>\r\n"))

				// Read incoming message data
				var data []byte
				for {
					n, err := conn.Read(buf)
					if err != nil {
						log.Println(err)
						break
					}
					data = append(data, buf[:n]...)
					if strings.HasSuffix(string(data), "\r\n.\r\n") {
						break
					}
				}
				log.Debugf("Received message: \n%s\n", string(data))
				conn.Write([]byte("250 OK\r\n"))
			case "QUIT":
				if len(args) > 0 {
					log.Tracef("Received command %s %+v", command, args)
					conn.Write([]byte("501 Syntax: QUIT\r\n"))
					break
				}
				log.Tracef("Received command %s", command)
				conn.Write([]byte("221 Bye\r\n"))
				conn.Close()
				log.Debugf("Connection from %s closed\n", conn.RemoteAddr().String())
				return
			default:
				log.Debugf("Received unknown command %s with args: %s", command, strings.Join(args, ","))
				conn.Write([]byte("502 Command not implemented\r\n"))
		}
	}
}
