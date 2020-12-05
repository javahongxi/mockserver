package main

import (
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/javahongxi/mockserver/config"
	"github.com/javahongxi/mockserver/generator/city"
	"github.com/javahongxi/mockserver/generator/citylist"
	"github.com/javahongxi/mockserver/generator/profile"
	"github.com/javahongxi/mockserver/recommendation"
)

const templateSuggestion = "Please make sure working directory is the root of the repository, where we have go.mod/go.sum. Suggested command line: go run main.go"

func main() {
	profileTemplate, err := template.ParseFiles("generator/profile/profile_tmpl.html")
	if err != nil {
		log.Fatalf("Cannot create profile template: %v. %s", err, templateSuggestion)
	}
	profileGen := &profile.Generator{
		Tmpl:           profileTemplate,
		Recommendation: recommendation.Client{},
	}

	cityTemplate, err := template.ParseFiles("generator/city/city_tmpl.html")
	if err != nil {
		log.Fatalf("Cannot create city template: %v. %s", err, templateSuggestion)
	}
	cityGen := &city.Generator{
		Tmpl:       cityTemplate,
		ProfileGen: profileGen,
	}

	cityListTemplate, err := template.ParseFiles("generator/citylist/citylist_tmpl.html")
	if err != nil {
		log.Fatalf("Cannot create citylist template: %v. %s", err, templateSuggestion)
	}
	cityListGen := &citylist.Generator{
		Tmpl: cityListTemplate,
	}

	rand.Seed(time.Now().Unix())
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/static/index.html")
	})
	r.Static("/static", "static")
	r.GET("mock/www.zhenai.com/zhenghun", cityListGen.HandleRequest)
	r.GET("mock/www.zhenai.com/zhenghun/:city/:page", cityGen.HandleRequest)
	r.GET("mock/www.zhenai.com/zhenghun/:city", cityGen.HandleRequest)
	r.GET("mock/album.zhenai.com/u/:id", profileGen.HandleRequest)

	log.Fatal(r.Run(config.ListenAddress))
}
