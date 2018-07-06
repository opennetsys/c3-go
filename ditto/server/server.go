package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/c3systems/c3/ditto/util"
)

var listener net.Listener
var ipfsGateway = "http://127.0.0.1:9001"
var serverHost = "0.0.0.0:5000"

// Run ...
func Run() {
	if listener != nil {
		return
	}
	var gw string

	contentTypes := map[string]string{
		"manifestV2Schema":     "application/vnd.docker.distribution.manifest.v2+json",
		"manifestListV2Schema": "application/vnd.docker.distribution.manifest.list.v2+json",
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		uri := r.RequestURI
		fmt.Println(uri)

		if uri == "/v2/" {
			jsonstr := []byte(fmt.Sprintf(`{"what": "a registry", "gateway":%q, "handles": [%q, %q], "problematic": ["version 1 registries"], "project": "https://github.com/c3systems/c3"}`, gw, contentTypes["manifestListV2Schema"], contentTypes["manifestV2Schema"]))

			w.Header().Set("Docker-Distribution-API-Version", "registry/2.0")
			fmt.Fprintln(w, string(jsonstr))
			return
		}

		if len(uri) <= 1 {
			fmt.Fprintln(w, "invalid multihash")
			return
		}

		var suffix string
		if strings.HasSuffix(uri, "/latest") {
			// docker daemon requesting the manifest
			suffix = "-v1"
			accepts := r.Header["Accept"]
			for _, accept := range accepts {
				if accept == contentTypes["manifestV2Schema"] ||
					accept == contentTypes["manifestListV2Schema"] {
					suffix = "-v2"
					break
				}
			}
		}

		s := strings.Split(uri, "/")
		hash := util.IpfsifyHash(s[2])
		rest := strings.Join(s[3:], "/") // tag
		path := hash + "/" + rest

		// blob request
		location := ipfsGateway + "/ipfs/" + path

		if suffix != "" {
			// manifest request
			location = location + suffix
		}
		fmt.Printf("location %s", location)

		req, err := http.NewRequest("GET", location, nil)
		if err != nil {
			log.Fatal(err)
		}

		httpClient := http.Client{}
		resp, err := httpClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		//w.Header().Set("Location", location) // not required since we're fetching the content and proxying
		w.Header().Set("Docker-Distribution-API-Version", "registry/2.0")

		// if latest-v2 set header
		w.Header().Set("Content-Type", contentTypes["manifestV2Schema"])
		fmt.Fprintf(w, string(body))
	})

	var err error
	listener, err = net.Listen("tcp", serverHost)
	if err != nil {
		log.Fatal(err)
	}

	port := listener.Addr().(*net.TCPAddr).Port
	log.Println("PORT", port)

	log.Fatal(http.Serve(listener, nil))
}

// Close ...
func Close() {
	listener.Close()
}
