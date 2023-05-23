package main

import (
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
	"strings"
	"strconv"
	"github.com/kkyr/fig"
)

type config struct {
	Default		string		`fig:"default-language" default:"en"`
	Languages	[]string	`fig:"languages" default:"[\"en\"]"`
	RemoveRegion	bool		`fig:"remove-region"`
	Port		int		`fig:"port" default:"9000"`
	Address		string		`fig:"address" default:"localhost"`
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
	l, err := net.Listen("tcp", cfg.Address + ":" + strconv.Itoa(cfg.Port))
        if err != nil {
                log.Fatalln(err)
        }
	b := new(FastCGIServer)
        if err := fcgi.Serve(l, b); err != nil {
                log.Fatalln(err)
        }
}
