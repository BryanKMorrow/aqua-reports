package reports

import (
	"fmt"
	"github.com/BryanKMorrow/aqua-reports/src/system/aqua"
	"github.com/BryanKMorrow/aqua-reports/src/system/reports"
	"log"
	"net/http"
	"strconv"
	"time"
)

// AllStream - Stream results instead of JSON return
func AllStream(w http.ResponseWriter, r *http.Request) {
	defer Track(RunningTime("/reports/all/streams"))

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Server does not support Flusher!",
			http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	start := time.Now()
	csp := aqua.NewCSP()
	csp.ConnectCSP()

	pagesize := 100
	page := 1

	list, imageCount, repoCount, _ := csp.GetAllImages(pagesize, page)
	fmt.Fprintf(w, "Image Repository Count: %d - Total Image Count: %d - Pagesize (repos per query): %d\n", repoCount, imageCount, pagesize)
	count := 0
	if repoCount <= pagesize {
		for _, l := range list {
			for _, v := range l.Result {
				count++
				ir, exists := csp.GetImageRisk(v.Registry, v.Repository, v.Tag)
				if exists {
					vuln := csp.GetImageVulnerabilities(v.Registry, v.Repository, v.Tag)
					sens := csp.GetImageSensitive(v.Registry, v.Repository, v.Tag)
					malw := csp.GetImageMalware(v.Registry, v.Repository, v.Tag)
					_, path := reports.WriteHTMLReport(ir.Repository, ir.Tag, ir, vuln, malw, sens)
					url := fmt.Sprintf("http://%s/%s", r.Host, path)
					fmt.Fprintf(w, "Image scan report completed for %s -  %d of %d images (%d%%) - total time elapsed: %v\n",
						url, count, imageCount, count*100/imageCount, time.Since(start).Truncate(time.Millisecond))
					flusher.Flush()
				}
			}
		}
	} else {
		log.Printf("Repository Page: %d", page)
		for _, l := range list {
			for _, v := range l.Result {
				count++
				ir, exists := csp.GetImageRisk(v.Registry, v.Repository, v.Tag)
				if exists {
					vuln := csp.GetImageVulnerabilities(v.Registry, v.Repository, v.Tag)
					sens := csp.GetImageSensitive(v.Registry, v.Repository, v.Tag)
					malw := csp.GetImageMalware(v.Registry, v.Repository, v.Tag)
					_, path := reports.WriteHTMLReport(ir.Repository, ir.Tag, ir, vuln, malw, sens)
					url := fmt.Sprintf("http://%s/%s", r.Host, path)
					fmt.Fprintf(w, "Image scan report completed for %s -  %d of %d images (%d%%) - total time elapsed: %v\n",
						url, count, imageCount, count*100/imageCount, time.Since(start).Truncate(time.Millisecond))
					flusher.Flush()
				}
			}
		}

		currentCount := repoCount - pagesize
		log.Println("Outside for loop currentCount: " + strconv.Itoa(currentCount))
		if currentCount >= pagesize {
			page++
			log.Printf("Next Page: %d - Current Image Count: %d", page, currentCount)
		}
	}
	fmt.Fprintf(w, "Finished scan report creation for %d images\n", count)
}
