package test

import (
	"encoding/json"
	"github.com/MCPutro/golang-docker/model/web"
	"github.com/MCPutro/golang-docker/util"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"testing"
)

func TestWriteToResponseBody(t *testing.T) {

	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})

	dummyResponse := web.Response{
		Status:  200,
		Message: "success",
		Data:    "ok",
	}

	type args struct {
		c          *fiber.Ctx
		statusCode int
		message    string
		data       interface{}
	}

	type expectedResp struct {
		statusCode int
		message    string
	}

	tests := []struct {
		name         string
		args         args
		wantErr      bool
		expectedResp expectedResp
	}{
		// TODO: Add test cases.
		{
			name: "print response",
			args: args{
				c:          ctx,
				statusCode: fiber.StatusOK,
				message:    "success",
				data:       "ok",
			},
			wantErr: false,
			expectedResp: expectedResp{
				statusCode: 200,
				message:    "success",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := util.WriteToResponseBody(tt.args.c, tt.args.statusCode, tt.args.message, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("WriteToResponseBody() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, tt.args.c.Response().Header.StatusCode(), tt.expectedResp.statusCode)

			data, _ := json.Marshal(dummyResponse)
			assert.Equal(t, string(tt.args.c.Response().Body()), string(data))

		})
	}
}
