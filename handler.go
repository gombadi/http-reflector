package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type requestData struct {
	contentType   string
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
	Time          time.Time
}

// reflectHandler processes all requests and returns output in the requested format
func reflectHandler(w http.ResponseWriter, r *http.Request) {

	ipaddr, port := extractIP(r)

	rd := &requestData{
		IP:            ipaddr,
		Port:          port,
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

	switch {
	case strings.HasPrefix(rd.URLPath[1:], "all"):
		ob = writeAll(rd)
	default:
		ob = []byte(rd.IP + "\n")
	}

	if ob == nil {
		log.Printf("reflector: nil output buffer - sending internal server error\n")
		w.WriteHeader(500)
	} else {
		if len(rd.contentType) != 0 {
			w.Header().Set("Content-Type", rd.contentType)
		} else {
			w.Header().Set("Content-Type", "text/plain")
		}
		io.WriteString(w, string(ob))
	}

	log.Printf("reflector: %s %s %s\n",
		rd.IP,
		rd.Method,
		rd.RequestURI)
}

// extractIP extracts the ip & port from the http.Request.RemoteAddr field.
// This field is in different formats depending on ipv4/ipv6 and if the
// port info is available
func extractIP(r *http.Request) (ipaddr, port string) {

	// First check if there is a header X-Forwarded-For or similar
	for k, v := range r.Header {
		if ok := strings.Contains(strings.ToLower(k), "x-forwarded-for"); ok == true {
			// only return the first or only ip address
			return strings.Split(v[0], ",")[0], ""
		}
	}

	ipaddr, port, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", ""
	}
	return ipaddr, port
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

	switch {
	case strings.Contains(rd.URLPath, "/json"):
		return "json"
	case strings.Contains(rd.URLPath, "/xml"):
		return "xml"
	case strings.Contains(rd.URLPath, "/html"):
		return "html"
	}

	return "text"
}

// writeAll will return all requesrt information in the requested format
func writeAll(rd *requestData) []byte {

	switch selectOutput(rd) {
	case "json":

		b, err := json.MarshalIndent(rd, "", "\t")
		if err != nil {
			log.Printf("error with json.Marshal: %e", err)
		}
		rd.contentType = "application/json"
		return b

	case "xml":
		b, err := xml.MarshalIndent(rd, "", "\t")
		if err != nil {
			log.Printf("error with xml.Marshal: %e", err)
		}
		rd.contentType = "application/xml"
		return b

	default:

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
}

/*

 */
