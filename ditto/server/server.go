package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/c3systems/c3/ditto/util"
)

// Run ...
func Run() {
	var gw string

	contentTypes := map[string]string{
		"manifestV2Schema":     "application/vnd.docker.distribution.manifest.v2+json",
		"manifestListV2Schema": "application/vnd.docker.distribution.manifest.list.v2+json",
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		uri := r.RequestURI
		fmt.Println(uri)

		if uri == "/v2/" {
			jsonstr := []byte(fmt.Sprintf(`{"what": "a registry", "gateway":%q, "handles": [%q, %q], "project": "https://github.com/c3systems/c3"}`, gw, contentTypes["manifestListV2Schema"], contentTypes["manifestV2Schema"]))

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
		rest := strings.Join(s[3:], "/")
		path := hash + "/" + rest

		location := "http://localhost:9001/ipfs/" + path
		location = location + suffix
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

		//w.Header().Set("Location", location)
		w.Header().Set("Docker-Distribution-API-Version", "registry/2.0")

		// if latest-v2 set header
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Type", contentTypes["manifestV2Schema"])
		fmt.Fprintf(w, string(body))
	})

	listener, err := net.Listen("tcp", "0.0.0.0:5000")
	if err != nil {
		log.Fatal(err)
	}

	port := listener.Addr().(*net.TCPAddr).Port
	fmt.Println("PORT", port)

	log.Fatal(http.Serve(listener, nil))
}

func main() {
	Run()
}
