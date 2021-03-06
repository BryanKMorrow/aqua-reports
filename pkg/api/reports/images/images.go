package images

import (
	"encoding/json"
	"strconv"

	"log"
	"net/http"

	"github.com/BryanKMorrow/aqua-reports/pkg/api/reports"
	"github.com/gorilla/mux"
)

// Images sets a slice of images
type Images []Image

// AllHandler needs to handle the incoming request and execute the proper Image call
func AllHandler(w http.ResponseWriter, r *http.Request) {
	var images Images
	params := mux.Vars(r)
	queue := make(chan reports.Response)
	response := images.Get(params, queue)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Get - creates reports for all images in all registries
// Param: map[string]string - Map of request parameters
// Param: chan reports.Response - Channel that accepts the JSON response from each image
// Return: reports.Response - the Json response sent to the requester
func (il *Images) Get(params map[string]string, queue chan reports.Response) reports.Response {
	defer reports.Track(reports.RunningTime("images.Get"))
	var remaining, total, next int
	reports.UnescapeURLQuery(params)
	cli := reports.GetClient(il)
	all, remaining, next, total := cli.GetAllImages(0, 1000, nil, nil)
	log.Printf("Images: %d  Remaining: %d  Next Page: %d", total, remaining, next)
	for _, img := range all.Result {
		var image Image
		var p = make(map[string]string)
		p["registry"] = img.Registry
		p["image"] = img.Repository
		p["tag"] = img.Tag
		reports.UnescapeURLQuery(p)
		go image.Get(p, queue)
	}
	queueCount := 1
	for range queue {
		if queueCount == total {
			close(queue)
		}
		queueCount++
	}
	var response = reports.Response{
		Message: "Scan Report for all images: " + strconv.Itoa(total),
		URL:     "",
		Status:  "Write Successful",
	}
	return response
}
