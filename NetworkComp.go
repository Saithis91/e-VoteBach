package main

func InitSocket(IP string, port string, obj interface{}) {

	// Set port
	server.Port = port

	// Begin listening
	ln, err := net.Listen("tcp", IP+":"+port)
	if err != nil {
		panic(err)
	}

	// Close somehow
	defer ln.Close()

	// Log we're listening
	fmt.Println("Listening on IP and Port: " + ln.Addr().String())

	// While running - accept incoming connections
	for {

		// Accept
		conn, _ := ln.Accept() // Should do error checking here...
		obj.mutex.Lock()
		// Store connection
		obj.Serverconnections[conn.RemoteAddr().String()] = &conn
		obj.mutex.Unlock()
		// Handle connection
		go obj.HandleConnection(&conn)

	}
}