package download

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"
)

// Download all the files (might take several minutes).
func Download(dir string, timeout time.Duration, urlsOnly, skip bool, parallel, retries int) error {
	c := &http.Client{Timeout: timeout}
	if !urlsOnly {
		log.Output(2, "Preparing to download from the Federal Revenue official website…")
	}

	searches := []search{
		{
			"federal revenue",
			federalRevenue,
			federalRevenueSelector,
			federalRevenueExtension,
		},
		{
			"national treasure",
			nationalTreasure,
			nationalTreasureSelector,
			nationalTreasureExtension,
		},
	}
	fs, err := getFiles(c, searches, dir, skip)
	if err != nil {
		return err
	}

	if urlsOnly {
		urls := make([]string, 0, len(fs))
		for _, f := range fs {
			urls = append(urls, f.url)
		}
		sort.Strings(urls)
		for _, u := range urls {
			fmt.Println(u)
		}
		return nil
	}

	d, err := newDownloader(c, fs, 2, 4)
	if err != nil {
		return err
	}
	return d.downloadAll()
}
