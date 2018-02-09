package  easy_server
import(
	"fmt"
	"net"
	"runtime"
	"sync"
)

/*
type EasyServer and its functions
*/
type EasyServer struct{
	r receiver
	w worker
	num_of_workers int
	waitGroup * sync.WaitGroup
	num_of_listeners int
	tcpDataFuncPacketCh chan tcpDataFuncPacket
}

/*
create EasyServer object
*/
func NewServer(num int) *EasyServer {
    s := &EasyServer{
	tcpDataFuncPacketCh:        make(chan tcpDataFuncPacket),
		waitGroup: &sync.WaitGroup{},
		num_of_workers: num,
    }

    s.r = receiver{s.waitGroup,s.tcpDataFuncPacketCh}
    s.w = worker{s.waitGroup,s.tcpDataFuncPacketCh}

    return s
}

/*
listen on the specified port using tcp protocol
*/
func (server *EasyServer) AddTcpListener(port string,h * TcpDataHandlers){
	server.waitGroup.Add(1)
	server.num_of_listeners++
	go server.addTcpListener(port,h)
}

/*
create several workers
*/
func (server * EasyServer) CreateWorkers(){
	for i:=0;i<server.num_of_workers;i++ {
		server.waitGroup.Add(1)
		go server.w.handleTcpPacket(i)
	}
}

/*
print the internal information of EasyServer
*/
func (server * EasyServer) PrintServerInfo(){
	fmt.Println("number of listeners : ",server.num_of_listeners)
	fmt.Println("number of receivers : ",runtime.NumGoroutine()-server.num_of_listeners-server.num_of_workers)
	fmt.Println("number of workers : ",server.num_of_workers)
}

/*
wait all the works done
*/
func (server * EasyServer) Stop(){
	server.waitGroup.Wait()
	close(server.tcpDataFuncPacketCh)
}


func (server *EasyServer) addTcpListener(port string,h * TcpDataHandlers) {
	defer server.waitGroup.Done()
    ln, err := net.Listen("tcp",port)
    if err != nil {
        fmt.Println(err);
		panic("TCP can't listen on port "+port)
    }

	//bind this goroutine with a os.Thread
    runtime.LockOSThread()
    for {
        conn, err := ln.Accept()
        if err != nil {
            fmt.Println(err);
            continue
        }

		server.waitGroup.Add(1)
        go server.handleConnection(conn,h)
    }

}

func (server *EasyServer) handleConnection(conn net.Conn,h * TcpDataHandlers){
	defer server.waitGroup.Done()
	defer conn.Close()

	t := newTcpConnection(conn)
	server.waitGroup.Add(1)
	go server.r.splitPacket(t,h)

	for{
		select{
		case  <-t.closeCh:                    
			return
		case d:= <-t.dataCh:
                     if d!=nil {
			conn.Write(d)
                     }
		}
	}
}