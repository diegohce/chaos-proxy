package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func createProxies() *http.ServeMux {

	mux := http.NewServeMux()

	for path, hostconfig := range chaosConfig.Paths {
		log.Info().Println(path, "->", hostconfig.Host)

		url, err := url.Parse(hostconfig.Host)
		if err != nil {
			log.Error().Println(err, "parsing", hostconfig.Host)
			continue
		}

		p := httputil.NewSingleHostReverseProxy(url)
		p.ErrorHandler = errorHandler
		p.ModifyResponse = modifyResponse
		mux.Handle(path, p)
	}

	if url, err := url.Parse(chaosConfig.DefaultHost.Host); err != nil {
		log.Error().Println(err, "parsing", chaosConfig.DefaultHost.Host)
	} else {
		p := httputil.NewSingleHostReverseProxy(url)
		p.ErrorHandler = errorHandler
		p.ModifyResponse = modifyResponse
		mux.Handle("/", p)
	}

	return mux
}

func errorHandler(w http.ResponseWriter, r *http.Request, e error) {

	if e.Error() == "HUP" {
		hj, ok := w.(http.Hijacker)
		if !ok {
			log.Error().Println("Connection could not be hijacked")
			w.WriteHeader(400)
			return
		}

		conn, _, err := hj.Hijack()
		if err != nil {
			log.Error().Println("Error hijacking connection", err)
			w.WriteHeader(400)
			return
		}
		conn.Close()
		log.Info().Println("Connection closed for", r.URL.String())

	} else if e.Error() == "5xx" {
		if s, ok := e.(random5xx); ok {
			statusCode := s.status()
			w.WriteHeader(statusCode)

			log.Info().Println("Sending status", statusCode, "for", r.URL.String())
		} else {
			log.Error().Println("errorHandler: Cannot type-cast random5xx error")
			w.WriteHeader(http.StatusBadGateway)
		}

	} else if e.Error() == "TIMEOUT" {
		w.WriteHeader(http.StatusGatewayTimeout)

	} else {
		log.Error().Printf("http: proxy error: %v", e)
		w.WriteHeader(http.StatusBadGateway)
	}
}

func modifyResponse(res *http.Response) error {

	dice := rollDices()

	switch dice.kind() {
	case "HUP":
		return fmt.Errorf("HUP")

	case "5xx":
		if s, ok := dice.(random5xx); ok {
			return s
		}
		log.Error().Println("modifyResponse: Cannot type-cast random5xx error")

	case "DELAY":
		if r, ok := dice.(delay); ok {
			r.wait(res.Request.URL.String())
			return fmt.Errorf("TIMEOUT")
		}
		log.Error().Println("modifyResponse: Cannot type-cast delay error")

	}
	log.Info().Println("Passing through for", res.Request.URL.String())
	return nil
}
