package docs

import (
	"backend/config"
	"backend/utils/logs"
	"backend/utils/problems"
	"encoding/json"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"net/http"
)

const Path = "/docs/"
const OpenAPIPath = "/openapi.json/"

func CreateOpenAPISpec() *openapi3.T {
	doc := &openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:   "Lumina back-end API",
			Version: "1.0.0",
		},
	}

	return doc
}

func GetOpenAPISpec(port string) http.HandlerFunc {
	logs.Success(fmt.Sprintf("OpenAPI Spec available on http://localhost:%s%s", port, OpenAPIPath))

	return func(w http.ResponseWriter, r *http.Request) {
		doc := CreateOpenAPISpec()
		if err := json.NewEncoder(w).Encode(doc); err != nil {
			p := problems.Problem{
				Type:          problems.HandlerProblem,
				ServerMessage: fmt.Sprintf("while trying to get the openapi3 docs -> %v", err),
				ClientMessage: "An unexpected error has occurred while trying to fetch documentation.",
				Status:        http.StatusInternalServerError,
			}
			p.Handle(w, r)
			return
		}
	}
}

func GetSwagger(port string) http.HandlerFunc {
	logs.Success(fmt.Sprintf("Swagger running on http://localhost:%s%s", port, Path))

	return func(w http.ResponseWriter, r *http.Request) {
		fs := http.FileServer(http.Dir(config.SwaggerPath))
		http.StripPrefix(Path, fs).ServeHTTP(w, r)
	}
}
