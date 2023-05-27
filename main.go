package main

import (
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
	"strings"
	"strconv"
	"github.com/kkyr/fig"
	"os"
)

type config struct {
	Default		string		`fig:"default-language" default:"en"`
	Languages	[]string	`fig:"languages" default:"[\"en\"]"`
	RemoveRegion	bool		`fig:"remove-region"`
	Network		struct {
		Type	string	`fig:"type" default:"tcp"`
		Port	int	`fig:"port" default:"9000"`
		Address	string	`fig:"address" default:"localhost"`
		Unix	string  `fig:"unix" default:"/run/lang302.sock"`
	}
}
var cfg config

func load() error {
	return fig.Load(&cfg,
		fig.File("lang302.yaml"),
		fig.Dirs(".", "/etc/lang302", "/usr/local/etc/lang302"), 
	)
}

func contains(a []string, s string) bool {
	for _, value := range a {
		if value == s {
			return true
		}
	}
	return false
}

type FastCGIServer struct{}
func (s FastCGIServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	lang := cfg.Default
	langs := strings.Split(req.Header.Get("Accept-Language"), ",")
	for _, v := range langs {
		s := strings.Split(v, ";")
		if len(s) > 0 {
			v = s[0]
		}
		if cfg.RemoveRegion {
			s = strings.Split(v, "-")
			if len(s) > 0 {
				v = s[0]
			}
		}
		if contains(cfg.Languages, v) {
			lang = v
			break
		}
	}
	w.WriteHeader(302)
	w.Header().Set("Location", req.URL.String() + lang)
}

func main() {
	if err := load(); err != nil {
		log.Fatalln(err)
	}
	var listener net.Listener
	if cfg.Network.Type == "tcp" {
		addr := cfg.Network.Address + ":" +
					strconv.Itoa(cfg.Network.Port)
		l, err := net.Listen("tcp", addr)
		if err != nil {
			log.Fatalln(err)
		}
		listener = l
		log.Println("Listening on", addr)
	} else if cfg.Network.Type == "unix" {
		os.Remove(cfg.Network.Unix)
		unixAddr, err := net.ResolveUnixAddr("unix", cfg.Network.Unix)
		if err != nil {
			log.Fatalln(err)
		}
		l, err := net.ListenUnix("unix", unixAddr)
		if err != nil {
			log.Fatalln(err)
		}
		listener = l
		log.Println("Listening on unix:" + cfg.Network.Unix)
	} else {
		log.Fatalln("invalid network type", cfg.Network.Type)
	}
	b := new(FastCGIServer)
        if err := fcgi.Serve(listener, b); err != nil {
                log.Fatalln(err)
        }
}
