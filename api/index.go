package api

import (
	"github.com/gin-gonic/gin"
	"github.com/steeve/pulsar/xbmc"
)

func Index(c *gin.Context) {
	c.JSON(200, xbmc.NewView("", xbmc.ListItems{
		{Label: "Search Movies", Path: UrlForXBMC("/movies/search")},
		{Label: "Popular Movies", Path: UrlForXBMC("/movies/popular")},
		{Label: "Movies by Genre", Path: UrlForXBMC("/movies/genres")},

		{Label: "Search Shows", Path: UrlForXBMC("/shows/search")},
		{Label: "Popular Shows", Path: UrlForXBMC("/shows/popular")},

		{Label: "Search", Path: UrlForXBMC("/search")},
	}))
}