package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type requestData struct {
	IP            string
	Port          string
	Method        string
	RequestURI    string
	RequestURIlen int
	URLPath       string
	URLPathlen    int
	Protocol      string
	Host          string
	URLQuery      string
	URLQuerylen   int
	URLFragment   string
	Header        map[string][]string `xml:"-"`
	//Header map[string][]string
	Time time.Time
}

// reflectHandler processes all requests and returns output in the requested format
func reflectHandler(w http.ResponseWriter, r *http.Request) {

	rd := &requestData{
		IP:            strings.Split(r.RemoteAddr, ":")[0],
		Port:          strings.Split(r.RemoteAddr, ":")[1],
		Method:        r.Method,
		Protocol:      r.Proto,
		Host:          r.Host,
		RequestURI:    r.RequestURI,
		RequestURIlen: len(r.RequestURI),
		URLQuery:      r.URL.RawQuery,
		URLQuerylen:   len(r.URL.RawQuery),
		URLPath:       r.URL.Path,
		URLPathlen:    len(r.URL.Path),
		URLFragment:   r.URL.Fragment,
		Time:          time.Now().UTC(),
		Header:        r.Header,
	}

	var ob []byte

	switch rd.URLPath[1:] {
	case "ip":
		ob = []byte(rd.IP + "\n")
	case "all":
		switch selectOutput(rd) {
		case "json":
			ob = writeJson(rd)
			w.Header().Set("Content-Type", "application/json")
		case "xml":
			ob = writeXML(rd)
			w.Header().Set("Content-Type", "application/xml")
		default:
			ob = writeText(rd)
		}
	default:
		ob = []byte("Nothing to see here. Move along please\n")
	}

	if ob == nil {
		log.Printf("reflector: nil output buffer - sending internal server error\n")
		w.WriteHeader(500)
	} else {
		io.WriteString(w, string(ob))
	}

	log.Printf("reflector: %s %s %s\n",
		rd.IP,
		rd.Method,
		rd.RequestURI)
}

// selectOutput will return the requested output format
// Can be json, xml, html or the default text
func selectOutput(rd *requestData) string {

	switch {
	case strings.Contains(rd.URLQuery, "o=json"):
		return "json"
	case strings.Contains(rd.URLQuery, "o=xml"):
		return "xml"
	case strings.Contains(rd.URLQuery, "o=html"):
		return "html"
	}

	return "text"
}

// writeXML sends output in XML format
// xxxxFIXxxxx having some issues with maps being an unsupported type
func writeXML(rd *requestData) []byte {

	b, err := xml.MarshalIndent(rd, "", "\t")
	if err != nil {
		log.Printf("error with xml.Marshal: %v\n", err)
	}
	return b
}

// writeJson sends output in json format
func writeJson(rd *requestData) []byte {

	b, err := json.MarshalIndent(rd, "", "\t")
	if err != nil {
		log.Printf("error with json.Marshal: %e", err)
	}
	return b
}

// writeText sends output in text format
func writeText(rd *requestData) []byte {
	var b bytes.Buffer

	b.WriteString("Request:\n")
	b.WriteString("request.Time: " + rd.Time.String() + "\n")
	b.WriteString("request.IP: " + rd.IP + "\n")
	b.WriteString("request.Port: " + rd.Port + "\n")
	b.WriteString("request.Method: " + rd.Method + "\n")
	b.WriteString("request.Proto: " + rd.Protocol + "\n")
	b.WriteString("request.Host: " + rd.Host + "\n")
	b.WriteString("request.RequestURI.length: " + strconv.Itoa(rd.RequestURIlen) + "\n")
	b.WriteString("request.RequestURI: " + rd.RequestURI + "\n")
	b.WriteString("\nURL:\n")
	b.WriteString("url.Path.length: " + strconv.Itoa(rd.URLPathlen) + "\n")
	b.WriteString("url.Path: " + rd.URLPath + "\n")
	b.WriteString("url.RawQuery.length: " + strconv.Itoa(rd.URLQuerylen) + "\n")
	b.WriteString("url.RawQuery: " + rd.URLQuery + "\n")
	b.WriteString("url.Fragment: " + rd.URLFragment + "\n")

	b.WriteString("\nHeaders:\n")
	for k, v := range rd.Header {
		b.WriteString("header." + k + ": " + v[0] + "\n")
	}

	return b.Bytes()
}

/*

*/
