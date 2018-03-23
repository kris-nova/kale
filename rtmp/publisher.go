package rtmp

import (
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/kubicorn/kubicorn/pkg/logger"
	"github.com/nareix/joy4/av/avutil"
	"github.com/nareix/joy4/av/pubsub"
	"github.com/nareix/joy4/format"

	"github.com/nareix/joy4/format/rtmp"
)

// OBSServerOptions are the options for the server component of the router
type OBSServerOptions struct {

	// Raw key to use to secure the OBS stream
	RawKey string

	// Bind address to use for the server
	BindAddress string

	// The port number to use for the stream
	WebClientPort int

	// ForwardFunction is a function that will be called to route the OBS stream to
	ClientFuncs []clientFunc
}

// Initialize the channels for the stream
type obsChannel struct {
	queue *pubsub.Queue
}

// clientFunc is a function type we can use to register N forward functions for the OBS bouncer
type clientFunc func()

// NewObsOptions returns the default options for the struct
func NewObsOptions() *OBSServerOptions {
	return &OBSServerOptions{
		RawKey:        "",
		BindAddress:   "",
		WebClientPort: 8089,
		ClientFuncs:   []clientFunc{
		//webClientForwardFunction,
		},
	}
}

var (
	obsLock  = &sync.RWMutex{}
	channels = map[string]*obsChannel{}
)

// ListenAndServe is a wrapper function based on the original work of github.com/netroby/go-rtmp-server
// it will run the OBS server based on arbitrary OBS options that can be passed into the function
func ListenAndServe(o *OBSServerOptions) error {

	//sKey := o.RawKey
	addr := fmt.Sprintf("%s:%d", o.BindAddress, o.WebClientPort)

	server := &rtmp.Server{}

	// The client component
	server.HandlePlay = func(conn *rtmp.Conn) {
		obsLock.RLock()
		ch := channels[conn.URL.Path]
		obsLock.RUnlock()

		if ch != nil {
			cursor := ch.queue.Latest()
			avutil.CopyFile(conn, cursor)
		}
	}

	// The publisher component
	server.HandlePublish = func(conn *rtmp.Conn) {
		streams, _ := conn.Streams()
		obsLock.Lock()
		logger.Info("Request [%s]", conn.URL.RequestURI())
		logger.Info("Key [%s]", conn.URL.Query().Get("key"))
		streamKey := conn.URL.Query().Get("key")
		if streamKey != o.RawKey {
			logger.Critical("Invalid key for stream!")
			return
		}
		ch := channels[conn.URL.Path]
		if ch == nil {
			ch = &obsChannel{}
			ch.queue = pubsub.NewQueue()
			ch.queue.WriteHeader(streams)
			channels[conn.URL.Path] = ch
		} else {
			ch = nil
		}
		obsLock.Unlock()
		if ch == nil {
			return
		}
		avutil.CopyPackets(ch.queue, conn)
		obsLock.Lock()
		delete(channels, conn.URL.Path)
		obsLock.Unlock()
		ch.queue.Close()
	}

	// Register any client functions
	for _, ffunc := range o.ClientFuncs {
		ffunc()
	}

	logger.Info("Forwarding...")
	logger.Info("You can now access the stream on port rtmp://[hostname]:1735/live")
	go http.ListenAndServe(addr, nil)

	server.ListenAndServe()
	return nil
}

type writeFlusher struct {
	httpflusher http.Flusher
	io.Writer
}

func (self writeFlusher) Flush() error {
	self.httpflusher.Flush()
	return nil
}

func init() {
	format.RegisterAll()
}
