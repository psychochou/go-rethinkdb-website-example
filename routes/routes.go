package routes

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"regexp"
	"runtime"
	"strings"

	"github.com/flosch/pongo2"
	"github.com/valyala/fasthttp"
	r "gopkg.in/gorethink/gorethink.v3"
)

var (
	_methodGet              = []byte("GET")
	_methodPost             = []byte("POST")
	_routeHome              = []byte("/")
	_queryListSkip          = []byte{91, 55, 49, 44, 91, 91, 55, 48, 44, 91, 91, 51, 51, 44, 91, 91, 52, 49, 44, 91, 91, 49, 53, 44, 91, 91, 49, 52, 44, 91, 34, 112, 115, 121, 99, 104, 111, 34, 93, 93, 44, 34, 97, 114, 116, 105, 99, 108, 101, 115, 34, 93, 93, 44, 91, 55, 52, 44, 91, 34, 99, 114, 101, 97, 116, 101, 65, 116, 34, 93, 93, 93, 93, 44, 34, 105, 100, 34, 44, 34, 116, 105, 116, 108, 101, 34, 44, 34, 105, 109, 97, 103, 101, 34, 93, 93, 44}
	_queryListSkipSufix     = []byte{93, 93, 44, 49, 48, 93, 93}
	_queryListSkipTagFirst  = []byte{91, 55, 49, 44, 91, 91, 55, 48, 44, 91, 91, 51, 51, 44, 91, 91, 52, 49, 44, 91, 91, 55, 56, 44, 91, 91, 49, 53, 44, 91, 91, 49, 52, 44, 91, 34, 112, 115, 121, 99, 104, 111, 34, 93, 93, 44, 34, 97, 114, 116, 105, 99, 108, 101, 115, 34, 93, 93, 44, 34}
	_queryListSkipTagSecond = []byte{34, 93, 44, 123, 34, 105, 110, 100, 101, 120, 34, 58, 34, 116, 97, 103, 115, 34, 125, 93, 44, 91, 55, 52, 44, 91, 34, 99, 114, 101, 97, 116, 101, 65, 116, 34, 93, 93, 93, 93, 44, 34, 105, 100, 34, 44, 34, 116, 105, 116, 108, 101, 34, 44, 34, 105, 109, 97, 103, 101, 34, 93, 93, 44}
	_queryListSkipTagThird  = []byte{93, 93, 44, 49, 48, 93, 93}

	_routeRegexArticle        = regexp.MustCompile("^/a/[0-9a-zA-Z_\\-]+$")
	_routeRegexTag            = regexp.MustCompile("^/t/.+$")
	_routeRegexApiListSkip    = regexp.MustCompile("^/api/list/[0-9]+$")
	_routeRegexApiListSkipTag = []byte("/api/lt")

	_templates   map[string]*pongo2.Template
	_session     *r.Session
	_fileCss     string
	_fileJs      string
	_titleSuffix = "EverStore"
	_mapTags     = map[string]string{
		"health":       "健康",
		"encyclopedia": "百科",
		"foods":        "美食",
		"tools":        "工具",
	}
)

const (
	__templateDirectory = "_layout"
	__createAtabase     = "[64,[[69,[[2,[0]],[65,[[10,[0]],{\"created\":0},[57,[\"psycho\"]]]]]],[93,[[59,[]],\"psycho\"]]]]"
	__createTable       = "[64,[[69,[[2,[0]],[65,[[10,[0]],{\u0022created\u0022:0},[60,[[14,[\u0022psycho\u0022]],\u0022articles\u0022]]]\r\n]]],[93,[[62,[[14,[\u0022psycho\u0022]]]],\u0022articles\u0022]]]]"
	__createIndex       = "[64,[[69,[[2,[0]],[65,[[10,[0]],{\u0022created\u0022:0},[75,[[15,[[14,[\u0022psycho\u0022]],\u0022articles\u0022]],\u0022createAt\u0022]]]]]],[93,[[77,[[15,[[14,[\u0022psycho\u0022]],\u0022articles\u0022]]]],\u0022createAt\u0022]]]]"
	__createIndexTags   = "[64,[[69,[[2,[0]],[65,[[10,[0]],{\u0022created\u0022:0},[75,[[15,[[14,[\u0022psycho\u0022]],\u0022articles\u0022]],\u0022tags\u0022],{\u0022multi\u0022:true}]]]]],[93,[[77,[[15,[[14,[\u0022psycho\u0022]],\u0022articles\u0022]]]],\u0022tags\u0022]]]]"
	__watiIndexTags     = "[140,[[15,[[14,[\u0022psycho\u0022]],\u0022articles\u0022]],\u0022tags\u0022]]"
	__databaseAddress   = "localhost:28015"
	__fieldDatabase     = "psycho"
	__fieldTable        = "articles"
)

func generateListSkip(skip []byte) []byte {
	var buffer bytes.Buffer
	buffer.Write(_queryListSkip)
	buffer.Write(skip)
	buffer.Write(_queryListSkipSufix)
	return buffer.Bytes()
}
func generateListSkipTag(tag string, skip []byte) []byte {
	var buffer bytes.Buffer
	buffer.Write(_queryListSkipTagFirst)
	buffer.WriteString(tag)
	buffer.Write(_queryListSkipTagSecond)
	buffer.Write(skip)
	buffer.Write(_queryListSkipTagThird)
	return buffer.Bytes()
}
func createDatabase() {
	var err error
	_, err = r.RawQuery([]byte(__createAtabase)).Run(_session)
	if err != nil {
		log.Fatalln(err)
	}
	_, err = r.RawQuery([]byte(__createTable)).Run(_session)
	if err != nil {
		log.Fatalln(err)
	}
	_, err = r.RawQuery([]byte(__createIndex)).Run(_session)
	if err != nil {
		log.Fatalln(err)
	}
	cursor, err := r.DB("psycho").Table("articles").IndexWait("createAt").Pluck("ready").Run(_session)
	if err != nil {
		log.Fatalln(err)
	}
	result, err := cursor.Interface()
	fmt.Println(result.([]interface{})[0].(map[string]interface{})["ready"])

	_, err = r.RawQuery([]byte(__createIndexTags)).Run(_session)
	if err != nil {
		log.Fatalln(err)
	}
	_, err = r.RawQuery([]byte(__watiIndexTags)).Run(_session)
	if err != nil {
		log.Fatalln(err)
	}

}
func Prepare() {
	if runtime.GOOS == "windows" {
		_fileCss, _ = getFileVersion("public/stylesheets/main.css")
		_fileJs, _ = getFileVersion("public/javascripts/app.js")
	} else {
		_fileCss, _ = getFileVersion("/var/psycho/public/stylesheets/main.css")
		_fileJs, _ = getFileVersion("/var/psycho/public/javascripts/app.min.js")
	}
	//_fileJs, _ := getFileVersion("public/javascripts/main.js")

	var err error
	_session, err = r.Connect(r.ConnectOpts{
		Address: __databaseAddress,
	})

	if err != nil {
		log.Fatalln(err)
	}

	_templates = make(map[string]*pongo2.Template)
	files, err := ioutil.ReadDir(__templateDirectory)
	if err != nil {
		log.Fatalln(err)
	}

	for index := 0; index < len(files); index++ {
		fileName := files[index].Name()
		if files[index].IsDir() || strings.HasPrefix(fileName, "_") || !strings.HasSuffix(fileName, ".html") {
			continue
		}

		template := pongo2.Must(pongo2.FromFile(path.Join(__templateDirectory, fileName)))
		if err != nil {

			log.Fatalln(err)
		}
		_templates[fileName] = template
	}
}
func Handle(ctx *fasthttp.RequestCtx) {
	url := ctx.Path()
	if bytes.Compare(ctx.Method(), _methodGet) == 0 {
		if bytes.Compare(url, _routeHome) == 0 {
			getHome(ctx)
			return
		}
		if _routeRegexArticle.Match(url) {
			getArticle(ctx)
			return
		}
		if _routeRegexTag.Match(url) {
			getTag(ctx)
			return
		}
		if _routeRegexApiListSkip.Match(url) {
			apiListSkip(ctx)
			return
		}
		if bytes.Compare(url, _routeRegexApiListSkipTag) == 0 {
			apiListSkipTag(ctx)
			return
		}
		if r, _ := regexp.Match("\\.(?:css|js|woff2)$", url); r {
			getFile(ctx)
			return
		}
		sendStatusCode(ctx, 404)

	} else if bytes.Compare(ctx.Method(), _methodPost) == 0 {
		if bytes.HasPrefix(url, []byte("/psycho/")) {

			if bytes.HasSuffix(url, []byte("token")) {
				postCreateToken(ctx)
				return
			}
			if !postValidToken(ctx) {
				sendStatusCode(ctx, 401)
				return
			}

			if bytes.HasSuffix(url, []byte("update")) {
				postUpdate(ctx)
				return
			}
			if bytes.HasSuffix(url, []byte("insert")) {
				postInsert(ctx)
				return
			}
			if bytes.HasSuffix(url, []byte("delete")) {
				postDelete(ctx)
				return
			}
			if bytes.HasSuffix(url, []byte("importAll")) {
				postImportAll(ctx)
				return
			}
		}
	} else {
		sendStatusCode(ctx, 404)
	}
}

func sendStatusCode(ctx *fasthttp.RequestCtx, statusCode int) {
	ctx.SetStatusCode(statusCode)
	ctx.SetConnectionClose()
}

func sendString(ctx *fasthttp.RequestCtx, str string) {
	ctx.Response.Header.Add("Content-Type", "text/html; charset=utf-8")
	_, err := ctx.WriteString(str)
	if err != nil {
		fmt.Println(err)
	}
}
func sendBuffer(ctx *fasthttp.RequestCtx, buffer []byte) {
	ctx.Response.Header.Add("Content-Type", "text/html; charset=utf-8")
	_, err := ctx.Write(buffer)
	if err != nil {
		fmt.Println(err)
	}
}
func sendJson(ctx *fasthttp.RequestCtx, b []byte) {
	ctx.Response.Header.Add("Content-Type", "application/json; charset=utf-8")
	_, err := ctx.Write(b)
	if err != nil {
		fmt.Println(err)
	}
}
