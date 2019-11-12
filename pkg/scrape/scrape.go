package scrape

import (
	"crypto/sha1"
	"encoding/hex"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/ProtonMail/go-appdir"
	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
	"github.com/xbapps/xbvr/pkg/common"
	"github.com/xbapps/xbvr/pkg/models"
)

var log = &common.Log
var appDir string
var cacheDir string

var siteCacheDir string
var sceneCacheDir string

var userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36"


func createCollector(domains ...string) *colly.Collector {
	c := colly.NewCollector(
		colly.AllowedDomains(domains...),
		colly.CacheDir(siteCacheDir),
		colly.UserAgent(userAgent),
	)

	c = createCallbacks(c)
	return c
}

func cloneCollector(c *colly.Collector) *colly.Collector {
	x := c.Clone()
	x = createCallbacks(x)
	return x
}

func createCallbacks(c *colly.Collector) *colly.Collector {
	const maxRetries = 15

	c.OnRequest(func(r *colly.Request) {
		attempt := r.Ctx.GetAny("attempt")

		if attempt == nil {
			r.Ctx.Put("attempt", 1)
		}

		log.Infoln("visiting", r.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		attempt := r.Ctx.GetAny("attempt").(int)

		if r.StatusCode == 429 {
			log.Errorln("Error:", r.StatusCode, err)

			if attempt <= maxRetries {
				unCache(r.Request.URL.String(), c.CacheDir)
				log.Errorln("Waiting 2 seconds before next request...")
				r.Ctx.Put("attempt", attempt+1)
				time.Sleep(2 * time.Second)
				r.Request.Retry()
			}
		}
	})

	return c
}

func logScrapeStart(id string, name string) {
	log.WithFields(logrus.Fields{
		"task":      "scraperProgress",
		"scraperID": id,
		"progress":  0,
		"started":   true,
		"completed": false,
	}).Infof("Starting %v scraper", name)
}

func logScrapeFinished(id string, name string) {
	log.WithFields(logrus.Fields{
		"task":      "scraperProgress",
		"scraperID": id,
		"progress":  0,
		"started":   false,
		"completed": true,
	}).Infof("Finished %v scraper", name)
}

func registerScraper(id string, name string, f models.ScraperFunc) {
	models.RegisterScraper(id, name, f)
}

func unCache(URL string, cacheDir string) {
	sum := sha1.Sum([]byte(URL))
	hash := hex.EncodeToString(sum[:])
	dir := path.Join(cacheDir, hash[:2])
	filename := path.Join(dir, hash)
	if err := os.Remove(filename); err != nil {
		log.Fatal(err)
	}
}

func updateSiteLastUpdate(id string) {
	var site models.Site
	err := site.GetIfExist(id)
	if err != nil {
		log.Error(err)
		return
	}
	site.LastUpdate = time.Now()
	site.Save()
}

func init() {
	appDir = appdir.New("xbvr").UserConfig()

	cacheDir = filepath.Join(appDir, "cache")

	siteCacheDir = filepath.Join(cacheDir, "site_cache")
	sceneCacheDir = filepath.Join(cacheDir, "scene_cache")

	_ = os.MkdirAll(appDir, os.ModePerm)
	_ = os.MkdirAll(cacheDir, os.ModePerm)

	_ = os.MkdirAll(siteCacheDir, os.ModePerm)
	_ = os.MkdirAll(sceneCacheDir, os.ModePerm)
}
