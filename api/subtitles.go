package api

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/elgatito/elementum/bittorrent"

	"github.com/elgatito/elementum/config"
	"github.com/elgatito/elementum/osdb"
	"github.com/elgatito/elementum/util/ip"
	"github.com/elgatito/elementum/xbmc"

	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
)

var subLog = logging.MustGetLogger("subtitles")

// SubtitlesIndex ...
func SubtitlesIndex(s *bittorrent.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		q := ctx.Request.URL.Query()

		xbmcHost, _ := xbmc.GetXBMCHostWithContext(ctx)

		playingFile := xbmcHost.PlayerGetPlayingFile()

		// Check if we are reading a file from Elementum
		if strings.HasPrefix(playingFile, ip.GetContextHTTPHost(ctx)) {
			playingFile = strings.Replace(playingFile, ip.GetContextHTTPHost(ctx)+"/files", config.Get().DownloadPath, 1)
			// not QueryUnescape in order to treat "+" as "+" in file name on FS
			playingFile, _ = url.PathUnescape(playingFile)
		}

		showID := 0
		if s.GetActivePlayer() != nil {
			showID = s.GetActivePlayer().Params().ShowID
		}
		payloads, preferredLanguage := osdb.GetPayloads(xbmcHost, q.Get("searchstring"), strings.Split(q.Get("languages"), ","), q.Get("preferredlanguage"), showID, playingFile)
		subLog.Infof("Subtitles payload: %#v", payloads)

		results, err := osdb.DoSearch(payloads, preferredLanguage)
		if err != nil {
			subLog.Errorf("Error searching subtitles: %s", err)
		}

		items := make(xbmc.ListItems, 0)

		for _, sub := range results {
			rating, _ := strconv.ParseFloat(sub.SubRating, 64)
			subLang := sub.LanguageName
			if subLang == "Brazilian" {
				subLang = "Portuguese (Brazil)"
			}
			item := &xbmc.ListItem{
				Label:     subLang,
				Label2:    sub.SubFileName,
				Icon:      strconv.Itoa(int((rating / 2) + 0.5)),
				Thumbnail: sub.ISO639,
				Path: URLQuery(URLForXBMC("/subtitle/%s", sub.IDSubtitleFile),
					"file", sub.SubFileName,
					"lang", sub.SubLanguageID,
					"fmt", sub.SubFormat,
					"dl", sub.SubDownloadLink),
				Properties: &xbmc.ListItemProperties{},
			}
			if sub.MatchedBy == "moviehash" {
				item.Properties.SubtitlesSync = trueType
			}
			if sub.SubHearingImpaired == "1" {
				item.Properties.SubtitlesHearingImpaired = trueType
			}
			items = append(items, item)
		}

		ctx.JSON(200, xbmc.NewView("", items))
	}
}

// SubtitleGet ...
func SubtitleGet(ctx *gin.Context) {
	q := ctx.Request.URL.Query()
	file := q.Get("file")
	dl := q.Get("dl")

	outFile, _, err := osdb.DoDownload(file, dl)
	if err != nil {
		subLog.Error(err)
		ctx.String(200, err.Error())
		return
	}

	ctx.JSON(200, xbmc.NewView("", xbmc.ListItems{
		{Label: file, Path: outFile.Name()},
	}))
}
