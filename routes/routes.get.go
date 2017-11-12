package routes

import (
	"bytes"
	"encoding/json"
	"log"
	"net/url"
	"regexp"

	"github.com/flosch/pongo2"

	"github.com/valyala/fasthttp"
	r "gopkg.in/gorethink/gorethink.v3"
)

func getHome(ctx *fasthttp.RequestCtx) {
	var buffer bytes.Buffer
	buffer.Write(_queryListSkip)
	buffer.WriteByte('0')
	buffer.Write(_queryListSkipSufix)
	c, err := r.RawQuery(buffer.Bytes()).Run(_session)
	if err != nil {
		log.Println(err)

		sendStatusCode(ctx, 500)
		return
	}
	var r []interface{}
	err = c.All(&r)
	if err != nil {
		log.Println(err)

		sendStatusCode(ctx, 500)
		return
	}

	str, err := _templates["index.html"].ExecuteBytes(pongo2.Context{
		"type":           0,
		"title":          _titleSuffix,
		"items":          r,
		"fileCss":        _fileCss,
		"fileJavaScript": _fileJs,
	})
	if err != nil {
		log.Println(err)
		sendStatusCode(ctx, 500)
		return
	}
	sendBuffer(ctx, str)
}
func generateBreadcrumb(tagGenerice interface{}) string {
	var buffer bytes.Buffer
	tagArray, ok := tagGenerice.([]interface{})
	if ok {
		for index := 0; index < len(tagArray); index++ {
			tag := tagArray[index].(string)
			humanTag := _mapTags[tag]
			if len(humanTag) < 1 {
				humanTag = tag
			}
			if index == 0 {
				buffer.WriteString("<li class=\"breadcrumb-item\"><a class=\"breadcrumb-link\" href=\"/t/")
				buffer.WriteString(urlencode(tag))
				buffer.WriteString("\">")
				buffer.WriteString(humanTag)
				buffer.WriteString("</a></li>")
			} else {
				buffer.WriteString("<li class=\"breadcrumb-item\"><div class=\"breadcrumb-guillemet material-icons\"></div><a class=\"breadcrumb-link\" href=\"/t/")
				buffer.WriteString(urlencode(tag))
				buffer.WriteString("\">")
				buffer.WriteString(humanTag)
				buffer.WriteString("</a></li>")
			}
		}
	}
	return buffer.String()
}
func generateHumanTag(tag string) string {
	humanTag := _mapTags[tag]
	if len(humanTag) > 0 {
		return humanTag
	}
	return tag
}
func queryWithTag(tag string, skip int, limit int) ([]interface{}, error) {
	cursor, err := r.DB(__fieldDatabase).Table(__fieldTable).GetAll(tag).OptArgs(r.GetAllOpts{Index: "tags"}).OrderBy(r.Desc("createAt")).Skip(skip).Limit(limit).Run(_session)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var result []interface{}
	err = cursor.All(&result)

	if err != nil {
		log.Println(err)

		return nil, err
	}
	return result, err
}
func getTag(ctx *fasthttp.RequestCtx) {
	indexPosition := bytes.LastIndex(ctx.Path(), []byte("/"))
	if indexPosition == -1 {
		sendStatusCode(ctx, 402)
		return
	}
	tagEncoded := string(ctx.Path()[indexPosition+1:])
	tag, err := url.QueryUnescape(tagEncoded)

	queryResult, err := queryWithTag(tag, 0, 10)
	if err != nil {
		log.Println(err)
		sendStatusCode(ctx, 500)
		return
	}
	buffer, err := _templates["index.html"].ExecuteBytes(pongo2.Context{
		"type":           1,
		"title":          generateHumanTag(tag) + " - " + _titleSuffix,
		"items":          queryResult,
		"fileJavaScript": _fileJs,
		"fileCss":        _fileCss,
	})
	if err != nil {
		sendStatusCode(ctx, 500)
		return
	}
	sendBuffer(ctx, buffer)
}
func getArticle(ctx *fasthttp.RequestCtx) {
	var buffer bytes.Buffer
	buffer.Write(_queryListSkip)
	buffer.WriteByte('0')
	buffer.Write(_queryListSkipSufix)
	p := bytes.LastIndex(ctx.Path(), []byte("/"))
	if p < 0 {
		sendStatusCode(ctx, 402)
		return
	}
	id := string(ctx.Path()[p+1:])
	c, err := r.DB(__fieldDatabase).Table(__fieldTable).Get(id).Pluck("title", "content", "tags").Run(_session)
	if err != nil {
		sendStatusCode(ctx, 500)
		return
	}
	var r []map[string]interface{}
	err = c.All(&r)
	if err != nil {
		sendStatusCode(ctx, 500)
		return
	}

	str, err := _templates["article.html"].ExecuteBytes(pongo2.Context{
		"type":           2,
		"title":          r[0]["title"].(string) + " - " + _titleSuffix,
		"content":        r[0]["content"],
		"breadcrumb":     generateBreadcrumb(r[0]["tags"]),
		"fileJavaScript": _fileJs,
		"fileCss":        _fileCss,
	})
	if err != nil {
		sendStatusCode(ctx, 500)
		return
	}
	sendBuffer(ctx, str)
}
func apiListSkip(ctx *fasthttp.RequestCtx) {
	skipBytes, err := lastAfterSlash(ctx.Path())
	if err != nil {
		log.Println(err)
		sendStatusCode(ctx, 500)
		return
	}
	if ok, _ := regexp.Match("[0-9]+", skipBytes); !ok {
		sendStatusCode(ctx, 402)
		return
	}
	cursor, err := r.RawQuery(generateListSkip(skipBytes)).Run(_session)
	if err != nil {
		log.Println(err)
		sendStatusCode(ctx, 500)
		return
	}
	var result []interface{}
	err = cursor.All(&result)
	if err != nil {
		log.Println(err)
		sendStatusCode(ctx, 500)
		return
	}
	resultBytes, err := json.Marshal(&result)
	if err != nil {
		log.Println(err)
		sendStatusCode(ctx, 500)
		return
	}
	sendJson(ctx, resultBytes)
}
func apiListSkipTag(ctx *fasthttp.RequestCtx) {
	queryArgs := ctx.QueryArgs()
	tag, err := url.QueryUnescape(string(queryArgs.Peek("tag")))
	if err != nil {
		log.Println(err)
		sendStatusCode(ctx, 500)
		return
	}
	skip := queryArgs.Peek("skip")

	if len(skip) == 0 {
		sendStatusCode(ctx, 402)
		return
	}

	rawQuery := generateListSkipTag(tag, skip)

	cursor, err := r.RawQuery(rawQuery).Run(_session)

	if err != nil {
		log.Println(err)
		sendStatusCode(ctx, 500)
		return
	}
	var result []interface{}
	err = cursor.All(&result)
	if err != nil {
		log.Println(err)
		sendStatusCode(ctx, 500)
		return
	}
	resultBytes, err := json.Marshal(&result)
	if err != nil {
		log.Println(err)
		sendStatusCode(ctx, 500)
		return
	}
	sendJson(ctx, resultBytes)

}
