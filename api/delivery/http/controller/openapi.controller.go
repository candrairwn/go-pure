package controller

import (
	"embed"
	"net/http"
)

type OpenApiController struct {
}

func NewOpenApiController() *OpenApiController {
	return &OpenApiController{}
}

var Openapi []byte
var SwaggerUI embed.FS

func (c OpenApiController) HandleGetServeOpenapi(version string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(Openapi)
	})
}

func (c OpenApiController) HandleGetSwaggerUI() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/docs/" {
			http.Redirect(w, r, "/api/docs/swagger", http.StatusFound)
			return
		}

		http.StripPrefix("/api/docs", http.FileServerFS(SwaggerUI)).ServeHTTP(w, r)
	})
}
