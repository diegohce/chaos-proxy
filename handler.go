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

	} else if e.Error() == "5xx" {
		if s, ok := e.(random5xx); ok {
			w.WriteHeader(s.status())
		} else {
			log.Error().Println("errorHandler: Cannot type-cast random5xx error")
		}

	} else {
		log.Error().Printf("http: proxy error: %v", e)
		w.WriteHeader(http.StatusBadGateway)
	}
}

func modifyResponse(res *http.Response) error {

	dice := rollDices()

	switch(dice.kind()) {
	case "HUP":
		return fmt.Errorf("HUP")

	case "5xx":
		if s, ok := dice.(random5xx); ok {
			return s
		} else {
			log.Error().Println("modifiyResponse: Cannot type-cast random5xx error")
		}
	case "DELAY":
		if r, ok := dice.(delay); ok {
			r.wait()
			return nil
		}
	}
	return nil
}

