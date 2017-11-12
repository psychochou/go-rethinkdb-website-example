package routes

import (
	"bytes"
	"encoding/json"
	"log"
	"time"

	"../jwt"
	"../shortid"
	"github.com/valyala/fasthttp"
	r "gopkg.in/gorethink/gorethink.v3"
)

func postCreateToken(ctx *fasthttp.RequestCtx) {
	body := ctx.PostBody()
	var result interface{}
	err := json.Unmarshal(body, &result)
	if err != nil {
		log.Println(err)
		sendStatusCode(ctx, 500)
		return
	}

	m, ok := result.(map[string]interface{})

	if !ok {
		sendStatusCode(ctx, 402)
		return
	}
	p := m["password"]
	if p != "ju3bInCL74qy" {
		sendStatusCode(ctx, 402)
		return
	}
	alogrithms := jwt.HmacSha256("LwWPWbxEUDlR38FJ_d2ilRTU")
	claims := jwt.NewClaim()
	claims.Set("Role", "Admin")
	claims.SetTime("exp", time.Now().Add(time.Duration(1)*time.Hour))
	str, err := alogrithms.Encode(claims)
	if err != nil {
		log.Println(err)
		sendStatusCode(ctx, 500)
		return
	}

	ctx.Response.Header.Add("Authorization", "Bearer "+str)
	sendStatusCode(ctx, 200)
}

func postDelete(ctx *fasthttp.RequestCtx) {

	b := ctx.PostBody()
	var result interface{}
	err := json.Unmarshal(b, &result)
	if err != nil {
		log.Println(err)
		sendStatusCode(ctx, 500)
		return
	}

	rm, ok := result.(map[string]interface{})
	if !ok {
		sendStatusCode(ctx, 402)
		return
	}
	idGeneric := rm["id"]
	id, ok := idGeneric.(string)
	if !ok {
		sendStatusCode(ctx, 402)
		return
	}
	if len(id) < 1 {
		sendStatusCode(ctx, 402)

		return
	}

	c, err := r.DB(__fieldDatabase).Table(__fieldTable).Get(id).Delete().Run(_session)
	if err != nil {
		log.Println(err)
		sendStatusCode(ctx, 500)
		return
	}

	i, err := c.Interface()
	if err != nil {
		log.Println(err)
		sendStatusCode(ctx, 500)
		return
	}
	b, err = json.Marshal(i)
	if err != nil {
		log.Println(err)
		sendStatusCode(ctx, 500)
		return
	}
	sendJson(ctx, b)

}

func postImportAll(ctx *fasthttp.RequestCtx) {
	buffer := ctx.PostBody()
	var result interface{}
	err := json.Unmarshal(buffer, &result)
	if err != nil {
		log.Println(err)
		sendStatusCode(ctx, 500)
		return
	}
	cursor, err := r.DB("psycho").Table("articles").Insert(result).Run(_session)

	if err != nil {
		log.Println(err)
		sendStatusCode(ctx, 500)
		return
	}

	r, err := cursor.Interface()
	if err != nil {
		log.Println(err)
		sendStatusCode(ctx, 500)
		return
	}

	b, err := json.Marshal(r)
	if err != nil {
		log.Println(err)
		sendStatusCode(ctx, 500)
		return
	}
	sendBuffer(ctx, b)
}

func postInsert(ctx *fasthttp.RequestCtx) {

	b := ctx.PostBody()
	var result interface{}
	err := json.Unmarshal(b, &result)
	if err != nil {
		log.Println(err)
		sendStatusCode(ctx, 500)
		return
	}
	m := make(map[string]interface{})

	rm, ok := result.(map[string]interface{})
	if !ok {
		sendStatusCode(ctx, 402)
		return
	}
	id, err := shortid.Generate()
	if err != nil {
		log.Println(err)
		sendStatusCode(ctx, 500)
		return
	}
	m["id"] = id
	m["title"] = rm["title"]
	m["content"] = rm["content"]
	m["image"] = rm["image"]
	m["tags"] = rm["tags"]
	m["createAt"] = makeTimestamp()
	m["updateAt"] = makeTimestamp()
	c, err := r.DB(__fieldDatabase).Table(__fieldTable).Insert(m).Run(_session)
	if err != nil {
		log.Println(err)
		sendStatusCode(ctx, 500)
		return
	}

	i, err := c.Interface()
	if err != nil {
		log.Println(err)
		sendStatusCode(ctx, 500)
		return
	}
	b, err = json.Marshal(i)
	if err != nil {
		log.Println(err)
		sendStatusCode(ctx, 500)
		return
	}
	sendJson(ctx, b)

}
func postUpdate(ctx *fasthttp.RequestCtx) {

	b := ctx.PostBody()

	// Parse JSON
	var result interface{}
	err := json.Unmarshal(b, &result)
	if err != nil {
		log.Println(err)
		sendStatusCode(ctx, 500)
		return
	}

	// Parse Document
	m := make(map[string]interface{})

	rm, ok := result.(map[string]interface{})
	if !ok {
		sendStatusCode(ctx, 402)
		return
	}
	id, ok := rm["id"].(string)
	if !ok && len(id) > 0 {
		sendStatusCode(ctx, 402)
		return
	}
	m["title"] = rm["title"]
	m["content"] = rm["content"]
	m["image"] = rm["image"]
	m["tags"] = rm["tags"]
	m["updateAt"] = makeTimestamp()

	// Update Document
	c, err := r.DB(__fieldDatabase).Table(__fieldTable).Get(id).Update(m).Run(_session)
	if err != nil {
		log.Println(err)
		sendStatusCode(ctx, 500)
		return
	}

	i, err := c.Interface()
	if err != nil {
		log.Println(err)
		sendStatusCode(ctx, 500)
		return
	}
	b, err = json.Marshal(i)
	if err != nil {
		log.Println(err)
		sendStatusCode(ctx, 500)
		return
	}
	sendJson(ctx, b)

}
func postValidToken(ctx *fasthttp.RequestCtx) bool {
	b := ctx.Request.Header.Peek("Authorization")
	bs := bytes.SplitN(b, []byte(" "), 2)
	alogrithms := jwt.HmacSha256("LwWPWbxEUDlR38FJ_d2ilRTU")
	err := alogrithms.Validate(string(bs[1]))
	if err != nil {
		log.Println(err, string(b))

		return false
	}
	return true
}
