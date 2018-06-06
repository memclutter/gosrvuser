package main

import (
	"time"

	"bytes"

	"github.com/asaskevich/govalidator"
	"github.com/valyala/fasthttp"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

// Handle GET /status. This endpoint need for health check.
func HandleStatus(ctx *fasthttp.RequestCtx) {
	data := ResponseDataStatus{}
	data.Time = time.Now().UTC()

	data.Health.Db = ctx.UserValue("mongodb.db") != nil
	data.Health.Amqp = ctx.UserValue("amqp.ch") != nil

	ctx.Response.SetStatusCode(fasthttp.StatusOK)
	ctx.SetUserValue("response.data", data)
}

// Handle POST /sign-up. This endpoint need for create new user.
func HandleSignUp(ctx *fasthttp.RequestCtx) {
	signUp, err := NewRequestSignUp(bytes.NewReader(ctx.Request.Body()))
	if err != nil {
		ctx.Response.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	if isValid, errs := signUp.Validate(); !isValid {
		ctx.Response.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetUserValue("response.data", errs)
		return
	}

	// TODO: encrypt password
	// TODO: save in database
	// TODO: notify message bus
	// TODO: set nil password in response

	ctx.Response.SetStatusCode(fasthttp.StatusOK)
	ctx.SetUserValue("response.data", signUp)
}

// Handle not allowed http method
func HandleMethodNotAllowed(ctx *fasthttp.RequestCtx) {
	ctx.Response.SetStatusCode(fasthttp.StatusMethodNotAllowed)
}

// Handle not found
func HandleNotFound(ctx *fasthttp.RequestCtx) {
	ctx.Response.SetStatusCode(fasthttp.StatusNotFound)
}
