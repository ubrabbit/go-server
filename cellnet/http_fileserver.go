package cellnet

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	_ "github.com/davyxu/cellnet/peer/http"
	"github.com/davyxu/cellnet/proc"
	_ "github.com/davyxu/cellnet/proc/http"
)

func StartFileServer(address string, shareDir string) {
	queue := cellnet.NewEventQueue()

	p := peer.NewGenericPeer("http.Acceptor", "httpfile", address, nil).(cellnet.HTTPAcceptor)
	p.SetFileServe(".", shareDir)
	proc.BindProcessorHandler(p, "http", nil)
	p.Start()
	queue.StartLoop()
	queue.Wait()
}
