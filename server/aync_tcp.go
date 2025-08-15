package server

import (
	"fmt"
	"log"
	"net"
	"syscall"
)

// FDConn wraps a file descriptor to implement io.ReadWriter
type FDConn struct {
	fd int
}

func (f *FDConn) Read(p []byte) (int, error) {
	n, err := syscall.Read(f.fd, p)
	if err != nil {
		// Convert syscall errors to more standard errors
		if err == syscall.ECONNRESET {
			return 0, fmt.Errorf("connection reset by peer")
		}
		if err == syscall.EPIPE {
			return 0, fmt.Errorf("broken pipe")
		}
		return n, err
	}
	if n == 0 {
		// EOF - client closed connection
		return 0, fmt.Errorf("client closed connection")
	}
	return n, nil
}

func (f *FDConn) Write(p []byte) (int, error) {
	n, err := syscall.Write(f.fd, p)
	if err != nil {
		if err == syscall.ECONNRESET {
			return 0, fmt.Errorf("connection reset by peer")
		}
		if err == syscall.EPIPE {
			return 0, fmt.Errorf("broken pipe")
		}
		return n, err
	}
	return n, nil
}

func RunAsyncTCPServer(config Config) error {
	log.Println("Starting Async TCP server on", config.Host, config.Port)

	//maximum clients to be accepted
	max_clients := 20000

	var con_clients int = 0

	//create a socket
	serverFD, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		return err
	}
	defer syscall.Close(serverFD) // The defer always make sure to run when function return

	// Set SO_REUSEADDR to avoid "address already in use" errors
	err = syscall.SetsockoptInt(serverFD, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	if err != nil {
		return err
	}

	//set the socket to Nonblocking
	err = syscall.SetNonblock(serverFD, true)
	if err != nil {
		return err
	}

	// Bind the socket to IP and port
	// Parse the IP address from config
	var addr [4]byte
	if config.Host == "0.0.0.0" {
		addr = [4]byte{0, 0, 0, 0} // Listen on all interfaces
	} else if config.Host == "127.0.0.1" {
		addr = [4]byte{127, 0, 0, 1} // localhost
	} else {
		// Parse IP address
		ip := net.ParseIP(config.Host)
		if ip == nil {
			return fmt.Errorf("invalid IP address: %s", config.Host)
		}
		ipv4 := ip.To4()
		if ipv4 == nil {
			return fmt.Errorf("IPv6 not supported: %s", config.Host)
		}
		addr = [4]byte{ipv4[0], ipv4[1], ipv4[2], ipv4[3]}
	}

	err = syscall.Bind(serverFD, &syscall.SockaddrInet4{
		Port: config.Port,
		Addr: addr,
	})
	if err != nil {
		return err
	}

	//listen on this socket
	err = syscall.Listen(serverFD, max_clients)
	if err != nil {
		return err
	}

	//Async IO
	//create a EPOLL instance
	epollFD, err := syscall.EpollCreate1(0)
	if err != nil {
		return err
	}
	defer syscall.Close(epollFD)

	// need to add the vents to monitor.
	// At this moment we have only server socket which we will monitor to accept any client
	// and once any client connect we will add to the list

	/* https://man7.org/linux/man-pages/man2/epoll_ctl.2.html
		int epoll_ctl(int epfd, int op, int fd,
	                     struct epoll_event *_Nullable event);
		struct epoll_event {
	           uint32_t      events;  /* Epoll events
	           epoll_data_t  data;    /* User data variable
	       };
		   https://man7.org/linux/man-pages/man3/epoll_event.3type.html

	*/
	var socketServerEvent syscall.EpollEvent = syscall.EpollEvent{
		Events: syscall.EPOLLIN,
		Fd:     int32(serverFD),
	}
	err = syscall.EpollCtl(epollFD, syscall.EPOLL_CTL_ADD, serverFD, &socketServerEvent)
	if err != nil {
		return err
	}

	/* creting events for EpollWait to hold the object */
	var events []syscall.EpollEvent = make([]syscall.EpollEvent, max_clients)

	/* Run the loop
	It will accept the client and add the client to the epoll list */
	for {
		/* check if any FD is ready for IO */
		nevents, e := syscall.EpollWait(epollFD, events, -1)
		if e != nil {
			log.Printf("EpollWait error: %v\n", e)
			continue
		}
		for i := 0; i < nevents; i++ {
			//if the IO means for server socket , it is a new client connection
			if int(events[i].Fd) == serverFD {
				fd, _, err := syscall.Accept(serverFD)
				if err != nil {
					log.Println("err", err)
					continue
				}
				con_clients++
				// Extract client IP and port from sockaddr
				// Uncomment below if you want to log client connections:
				/*
				   if sockAddr, ok := addr.(*syscall.SockaddrInet4); ok {
				   	clientIP := fmt.Sprintf("%d.%d.%d.%d:%d",
				   		sockAddr.Addr[0], sockAddr.Addr[1], sockAddr.Addr[2], sockAddr.Addr[3],
				   		sockAddr.Port)
				   	log.Printf("Client connected: %s, concurrent clients: %d\n", clientIP, con_clients)
				   } else {
				   	//log.Printf("Client connected (unknown address), concurrent clients: %d\n", con_clients)
				   }
				*/

				syscall.SetNonblock(fd, true) // Fix: set client fd to non-blocking
				//add this fd to be monitored for IO
				var socketClientEvent syscall.EpollEvent = syscall.EpollEvent{
					Events: syscall.EPOLLIN | syscall.EPOLLHUP | syscall.EPOLLERR,
					Fd:     int32(fd),
				}

				err = syscall.EpollCtl(epollFD, syscall.EPOLL_CTL_ADD, fd, &socketClientEvent)
				if err != nil {
					log.Printf("Error adding client fd %d to epoll: %v\n", fd, err)
					syscall.Close(fd)
					con_clients--
				}
			} else {
				/* if here means IO from an existing client */
				clientFD := int(events[i].Fd)

				// Check for error or hangup events
				if events[i].Events&(syscall.EPOLLHUP|syscall.EPOLLERR) != 0 {
					con_clients--
					syscall.EpollCtl(epollFD, syscall.EPOLL_CTL_DEL, clientFD, nil)
					syscall.Close(clientFD)
					continue
				}

				// Create wrapper to use with ReadCommand
				conn := &FDConn{fd: clientFD}

				command, err := ReadCommand(conn)
				if err != nil {
					// Check if it's a non-blocking "would block" error
					if err == syscall.EAGAIN || err == syscall.EWOULDBLOCK {
						continue // No data available right now
					}
					// Client disconnected or other error
					//log.Printf("Client disconnected (fd: %d), error: %v, concurrent clients: %d\n", clientFD, err, con_clients-1)
					con_clients--
					syscall.EpollCtl(epollFD, syscall.EPOLL_CTL_DEL, clientFD, nil)
					syscall.Close(clientFD)
					continue
				}

				if command == nil {
					// No data read, but no error - shouldn't happen with epoll
					log.Printf("No command read from fd: %d\n", clientFD)
					continue
				}

				err = Respond(conn, command)
				if err != nil {
					log.Printf("Error responding (fd: %d): %v, concurrent clients: %d\n", clientFD, err, con_clients-1)
					con_clients--
					syscall.EpollCtl(epollFD, syscall.EPOLL_CTL_DEL, clientFD, nil)
					syscall.Close(clientFD)
				}
			}
		}

	}
}
