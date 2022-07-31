package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/planet-i/goexample1/cacheExample/cache"
)

type Server struct {
	cache.Cache
}

//为两类HTTP端点的请求，定义两个Handler结构并分别实现它们的ServerHTTP方法
func (s *Server) Listen() {
	http.Handle("/cache/", s.cacheHandler())
	http.Handle("/status/", s.statusHandler())
	http.ListenAndServe(":12345", nil)
}
func New(c cache.Cache) *Server {
	return &Server{c}
}

//cacheHandler 内嵌了Server结构 Server结构内嵌了Cache接口
type cacheHandler struct {
	*Server
}

//实例化cacheHandler
func (s *Server) cacheHandler() http.Handler {
	return &cacheHandler{s}
}
func (h *cacheHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := strings.Split(r.URL.EscapedPath(), "/")[2]
	if len(key) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	m := r.Method
	if m == http.MethodPut {
		b, _ := ioutil.ReadAll(r.Body)
		if len(b) != 0 {
			e := h.Set(key, b)
			if e != nil {
				log.Println(e)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
		return
	}
	if m == http.MethodGet {
		b, e := h.Get(key)
		if e != nil {
			log.Println(e)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(b) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Write(b)
		return
	}
	if m == http.MethodDelete {
		e := h.Del(key)
		if e != nil {
			log.Println(e)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

//statusHandler 内嵌了Server结构 Server结构内嵌了Cache接口
type statusHandler struct {
	*Server
}

//实例化statusHandler
func (s *Server) statusHandler() http.Handler {
	return &statusHandler{s}
}
func (h *statusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	b, e := json.Marshal(h.GetStat())
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(b)
}
