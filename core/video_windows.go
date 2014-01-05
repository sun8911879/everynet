package core

import (
	"regexp"
	"strings"
)

var video_static_youku, _ = regexp.Compile(`^http:\/\/static\.youku\.com\/.*?q?(player|loader)(_[^.]+)?\.swf`)

//http://static.youku.com/v1.0.0393/v/swf/loader.swf

func (Tcp *Request) Video() {
	switch Tcp.Source {
	case `static.youku.com`:
		if video_static_youku.MatchString(Tcp.Path) == true {
			Tcp.Host = "sunloufile.qiniudn.com"
			index := strings.Index(Tcp.Path, "?")
			if index != -1 {
				Tcp.GET = "/youku_player.swf" + Tcp.Path[index:]
			} else {
				Tcp.GET = "/youku_player.swf"
			}
		}
		break
	}
}
