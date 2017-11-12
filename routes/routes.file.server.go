package routes

import (
	"path"

	"github.com/valyala/fasthttp"
)

func getFile(ctx *fasthttp.RequestCtx) {
	ctx.SendFile(path.Join("public", string(ctx.Path())))
}
