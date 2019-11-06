package reports

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/BryanKMorrow/aqua-reports/src/system/aqua"
	"github.com/BryanKMorrow/aqua-reports/src/system/reports"
	"github.com/gorilla/mux"
)

// One - Get one image
func One(w http.ResponseWriter, r *http.Request) {
	defer Track(RunningTime("One"))
	log.Println("/reports/one route called")
	w.Header().Set("Content-Type", "application/json")
	var responseList []ImageResponse

	params := mux.Vars(r)
	encodedRegistry := params["registry"]
	registry, err := url.QueryUnescape(encodedRegistry)
	if err != nil {
		log.Fatal(err)
		return
	}
	encodedImage := params["image"]
	image, err := url.QueryUnescape(encodedImage)
	if err != nil {
		log.Fatal(err)
		return
	}
	encodedTag := params["tag"]
	tag, err := url.QueryUnescape(encodedTag)
	if err != nil {
		log.Fatal(err)
		return
	}

	csp := aqua.NewCSP()
	csp.ConnectCSP()

	ir := csp.GetImageRisk(registry, image, tag)
	vuln := csp.GetImageVulnerabilities(registry, image, tag)
	sens := csp.GetImageSensitive(registry, image, tag)
	malw := csp.GetImageMalware(registry, image, tag)
	// Write HTML
	resp := reports.WriteHTMLReport(image, tag, ir, vuln, malw, sens)

	var response = ImageResponse{image, tag, registry, resp}
	responseList = append(responseList, response)

	json.NewEncoder(w).Encode(responseList)
}