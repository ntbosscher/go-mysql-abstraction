package mysql

import (
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native" // Native engine
)

var isInit bool

// max pipeline buffer size for stream option
const bufSize = 20

// max size of the connection pool
const connectionPool = 20

// channel to retrieve a connection from the pool
var conn chan mysql.Conn

func prep(){
	isInit = true
	conn = make(chan mysql.Conn, connectionPool)

	// fill the connection pool
	go func () {
		for i := 0; i < connectionPool; i++ {
			conn <- newConnection()
		}
	}()
}

func newConnection() (mysql.Conn) {
	conn := mysql.New("tcp", "", "127.0.0.1:3306", "root", "Theuser1", "blue-giraffe")
	conn.Connect()
	return conn
}

// worker for returnConnection.
// if the connection is still connected, reinsert into the queue
// if not close the connection
func _returnConnection(cn mysql.Conn){
	if cn.IsConnected() {
		conn <- cn
	} else {
		cn.Close()

		// replace the connection with a new one
		conn <- newConnection()
	}
}

// a connection that's been used and is now idle
func returnConnection(cn mysql.Conn){
	// make this a non-blocking operation
	go _returnConnection(cn)
}

// retrieve a new connection
func getConnection() (mysql.Conn){
	if !isInit {
		prep()
	}

	return <-conn
}
