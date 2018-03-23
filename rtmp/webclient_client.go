package rtmp

import (
	"github.com/nareix/joy4/format/flv"

	"io"
	"net/http"

	"github.com/nareix/joy4/av/avutil"
)

func webClientForwardFunction() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		obsLock.RLock()
		ch := channels[r.URL.Path]
		obsLock.RUnlock()

		if ch != nil {
			w.Header().Set("Content-Type", "video/x-flv")
			w.Header().Set("Transfer-Encoding", "chunked")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.WriteHeader(200)
			flusher := w.(http.Flusher)
			flusher.Flush()

			muxer := flv.NewMuxerWriteFlusher(writeFlusher{httpflusher: flusher, Writer: w})
			cursor := ch.queue.Latest()

			// Hacking in here to reduce latency

			//copyFile(muxer, cursor, flusher)

			avutil.CopyFile(muxer, cursor)
		} else {
			//logger.Info("Request url: ", r.URL.Path)
			if r.URL.Path != "/" {
				http.NotFound(w, r)
			} else {
				homeHtml := `
				<!DOCTYPE html>
<html>
	<head>
		<title>Demo live</title>
		<style>
		body {
			margin:0;
			padding:0;
			background:#000;
		}

        video {
			position:absolute;
			width:100%;
			height:100%;
		}
		</style>
	</head>
	<body>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/flv.js/1.3.2/flv.min.js"></script>

        <video id="videoElement" controls autoplay x5-video-player-type="h5" x5-video-player-fullscreen="true" playsinline webkit-playsinline>
            Your browser is too old which doesn't support HTML5 video.
        </video>
		<script>
if (flvjs.isSupported()) {
	var videoElement = document.getElementById('videoElement');
	var flvPlayer = flvjs.createPlayer({
		type: 'flv',
		url: '/live'
	});
	flvPlayer.attachMediaElement(videoElement);
	flvPlayer.load();
	flvPlayer.play();
}
		</script>
	</body>
</html>`
				io.WriteString(w, homeHtml)
			}
		}
	})
}
