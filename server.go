package main

import (
	"encoding/json"
	"fmt"
	"hash/crc64"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type tServer struct {
	httpServer http.Server
	ttl        int64
	threshold  uint
	cache      tCachedItems
}

func (s *tServer) handleReport(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		writeReason(res, http.StatusMethodNotAllowed, cMethodNotAllowedErr.Error())
		return
	}
	if contentType := req.Header.Get("Content-Type"); contentType != cContentTypeApplicationJson {
		writeReason(res, http.StatusUnsupportedMediaType, cWrongContentTypeErr.Error())
		return
	}
	req.Body = http.MaxBytesReader(nil, req.Body, cMaxBodySize)
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		writeReason(res, http.StatusRequestEntityTooLarge, err.Error())
		return
	}
	incomingRequestData := map[string]string{}
	if err := json.Unmarshal(body, &incomingRequestData); err != nil {
		writeReason(res, http.StatusUnprocessableEntity, err.Error())
		return
	}
	if len(incomingRequestData) != 1 || incomingRequestData["url"] == "" {
		writeReason(res, http.StatusUnprocessableEntity, cWrongPayloadObjectErr.Error())
		return
	}
	result := s.handlePayload(incomingRequestData["url"])
	if result == nil {
		res.WriteHeader(http.StatusNoContent)
		return
	}
	res.Write(result)
}

func (s *tServer) handlePayload(url string) []byte {
	crc64Int := crc64.Checksum([]byte(url), gCrc64Table)
	s.cache.RLock()
	val, ok := s.cache.data[crc64Int]
	s.cache.RUnlock()
	if !ok {
		val = &tItem{createdAt: time.Now().UnixNano(),
			cnt: 1}
		s.addToCache(crc64Int, val)
		return nil
	} else {
		now := time.Now().UnixNano()
		if now >= val.createdAt+s.ttl {
			s.cache.Lock()
			val.cnt = 1
			val.createdAt = now
			s.cache.Unlock()
			return nil
		}
		if val.cnt >= s.threshold {
			return cBlockedResult
		}
		s.cache.Lock()
		val.cnt++
		s.cache.Unlock()
	}
	return nil
}

func (s *tServer) addToCache(key uint64, item *tItem) {
	s.cache.Lock()
	s.cache.data[key] = item
	s.cache.Unlock()
}

func (s *tServer) cleanCache() {
	s.cache.Lock()
	defer s.cache.Unlock()
	now := time.Now().UnixNano()
	for key, val := range s.cache.data {
		if val.createdAt+s.ttl < now {
			delete(s.cache.data, key)
		}
	}
}

func (s *tServer) cleanCacheContinuesly() {
	for {
		s.cleanCache()
		time.Sleep(2 * time.Minute)
	}
}

func (s *tServer) start() error {
	httpServerMux := http.NewServeMux()
	s.httpServer = http.Server{
		Addr:    ":1489",
		Handler: httpServerMux,
	}
	httpServerMux.HandleFunc("/report", s.handleReport)
	fmt.Printf("%s\tstarting to listen on %s\n", time.Now().String(), s.httpServer.Addr)
	if err := s.httpServer.ListenAndServe(); err != nil {
		return fmt.Errorf("error while listenning and serving: %v", err)
	}
	return nil
}

func newServer(ttl int64, threshold uint) *tServer {
	cache := tCachedItems{new(sync.RWMutex), map[uint64]*tItem{}}
	s := &tServer{ttl: ttl * 1000000, /* convert to nanos for simplier calculations */
		threshold: threshold,
		cache:     cache}
	go s.cleanCacheContinuesly()
	return s
}

func writeReason(res http.ResponseWriter, status int, format string, a ...interface{}) {
	reason := fmt.Sprintf(format, a...)
	res.Header().Set("X-Reason", reason)
	res.WriteHeader(status)
}
