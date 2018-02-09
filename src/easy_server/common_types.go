package easy_server

type SplitError int
const (
	NoSplitError SplitError = iota
	LessDataSplitError
	OtherSplitError
)

type tcpDataFuncPacket struct{
	ops TcpConnectionOps
	bytes []byte
	handlers * TcpDataHandlers
}

type TcpDataHandlers struct{
	splitPacket func([]byte) (int,SplitError)
	handleFirstPacket func(TcpConnectionOps,[]byte)
	handleNoFirstPacket func(TcpConnectionOps,[]byte)
}

func NewTcpDataHandlers(s func([]byte) (int,SplitError),hf,ho func(TcpConnectionOps,[]byte) ) * TcpDataHandlers{
     t := &TcpDataHandlers{s,hf,ho}
     return t
}