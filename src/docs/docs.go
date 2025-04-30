package docs

import (
	"backend/config"
	"backend/utils/logs"
	"backend/utils/problems"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
)

const Path = "/docs/"
const OpenAPIPath = "/openapi.json/"

var StringType = openapi3.NewStringSchema().Type
var ObjectType = openapi3.NewObjectSchema().Type
var descriptions = []string{
	"User details retrieved successfully.",
	"User logged out successfully.",
	"User is logged in.",
	"User is not logged in.",
	"User registered successfully.",
	"Invalid request body.",
	"User with these credentials already exists.",
	"Unexpected server error.",
	"Login successful.",
	"Invalid credentials or user not verified.",
	"Email verified successfully.",
	"The verification link is invalid or has expired.",
	"The verification link has expired.",
	"User details updated successfully.",
	"Password change request initialized successfully.",
	"User with the provided email does not exist.",
	"Password changed successfully.",
	"The password change token is invalid or has expired.",
	"Password change token is valid.",
	"The password change token has expired.",
	"The password change token is invalid.",
}

func CreateOpenAPISpec() *openapi3.T {
	paths := openapi3.NewPaths()

	/* user/get-user endpoint */
	responses := openapi3.NewResponses()
	responses.Set("200", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &descriptions[0],
			Content: openapi3.Content{
				"application/json": &openapi3.MediaType{
					Schema: &openapi3.SchemaRef{
						Value: &openapi3.Schema{
							Type: ObjectType,
							Properties: map[string]*openapi3.SchemaRef{
								"username": {Value: &openapi3.Schema{Type: StringType}},
								"email":    {Value: &openapi3.Schema{Type: StringType}},
								"image":    {Value: &openapi3.Schema{Type: StringType}},
							},
						},
					},
				},
			},
		},
	})
	responses.Set("401", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &descriptions[3],
		},
	})
	responses.Set("500", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &descriptions[7],
		},
	})
	paths.Set("/user/get-user", &openapi3.PathItem{
		Get: &openapi3.Operation{
			Summary:   "Get user details.",
			Tags:      []string{"User"},
			Responses: responses,
		},
	})

	/* user/logged-in */
	responses = openapi3.NewResponses()
	responses.Set("200", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &descriptions[2],
		},
	})
	responses.Set("401", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &descriptions[3],
		},
	})
	responses.Set("500", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &descriptions[7],
		},
	})
	paths.Set("/user/logged-in", &openapi3.PathItem{
		Get: &openapi3.Operation{
			Summary:   "Check whether user is logged in.",
			Tags:      []string{"User"},
			Responses: responses,
		},
	})

	/* user/register endpoint */
	responses = openapi3.NewResponses()
	responses.Set("201", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &descriptions[4],
		},
	})
	responses.Set("400", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &descriptions[5],
		},
	})
	responses.Set("409", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &descriptions[6],
		},
	})
	responses.Set("500", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &descriptions[7],
		},
	})
	paths.Set("/user/register", &openapi3.PathItem{
		Post: &openapi3.Operation{
			Summary: "Register a new user.",
			Tags:    []string{"User"},
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Description: "User registration details",
					Required:    true,
					Content: openapi3.Content{
						"application/json": &openapi3.MediaType{
							Schema: &openapi3.SchemaRef{
								Value: &openapi3.Schema{
									Type: ObjectType,
									Properties: map[string]*openapi3.SchemaRef{
										"username": {Value: &openapi3.Schema{Type: StringType}},
										"email":    {Value: &openapi3.Schema{Type: StringType}},
										"password": {Value: &openapi3.Schema{Type: StringType}},
									},
									Required: []string{"username", "email", "password"},
								},
							},
						},
					},
				},
			},
			Responses: responses,
		},
	})

	/* user/login endpoint */
	responses = openapi3.NewResponses()
	responses.Set("200", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &descriptions[8],
		},
	})
	responses.Set("400", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &descriptions[5],
		},
	})
	responses.Set("401", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &descriptions[9],
		},
	})
	responses.Set("500", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &descriptions[7],
		},
	})
	paths.Set("/user/login", &openapi3.PathItem{
		Post: &openapi3.Operation{
			Summary: "Log into user's account.",
			Tags:    []string{"User"},
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Description: "Login credentials",
					Required:    true,
					Content: openapi3.Content{
						"application/json": &openapi3.MediaType{
							Schema: &openapi3.SchemaRef{
								Value: &openapi3.Schema{
									Type: ObjectType,
									Properties: map[string]*openapi3.SchemaRef{
										"email":    {Value: &openapi3.Schema{Type: StringType}},
										"password": {Value: &openapi3.Schema{Type: StringType}},
									},
									Required: []string{"email", "password"},
								},
							},
						},
					},
				},
			},
			Responses: responses,
		},
	})

	/* user/verify-email endpoint */
	responses = openapi3.NewResponses()
	responses.Set("204", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &descriptions[10],
		},
	})
	responses.Set("400", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &descriptions[5],
		},
	})
	responses.Set("404", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &descriptions[11],
		},
	})
	responses.Set("410", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &descriptions[12],
		},
	})
	responses.Set("500", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &descriptions[7],
		},
	})
	paths.Set("/user/verify-email", &openapi3.PathItem{
		Patch: &openapi3.Operation{
			Summary: "Verify user's email.",
			Tags:    []string{"User"},
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Description: "Email Verification Token",
					Required:    true,
					Content: openapi3.Content{
						"application/json": &openapi3.MediaType{
							Schema: &openapi3.SchemaRef{
								Value: &openapi3.Schema{
									Type: ObjectType,
									Properties: map[string]*openapi3.SchemaRef{
										"token": {Value: &openapi3.Schema{Type: StringType}},
									},
									Required: []string{"email"},
								},
							},
						},
					},
				},
			},
			Responses: responses,
		},
	})

	/* user/logout endpoint */
	responses = openapi3.NewResponses()
	responses.Set("200", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &descriptions[1],
		},
	})
	responses.Set("500", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &descriptions[7],
		},
	})
	paths.Set("/user/logout", &openapi3.PathItem{
		Delete: &openapi3.Operation{
			Summary:   "Logs user out of their account.",
			Tags:      []string{"User"},
			Responses: responses,
		},
	})

	/* user/modify-user endpoint */
	responses = openapi3.NewResponses()
	responses.Set("200", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &descriptions[13],
		},
	})
	responses.Set("400", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &descriptions[5],
		},
	})
	responses.Set("500", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &descriptions[7],
		},
	})
	paths.Set("/user/modify-user", &openapi3.PathItem{
		Patch: &openapi3.Operation{
			Summary: "Modify user details.",
			Tags:    []string{"User"},
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Description: "User details to be updated",
					Required:    true,
					Content: openapi3.Content{
						"application/json": &openapi3.MediaType{
							Schema: &openapi3.SchemaRef{
								Value: &openapi3.Schema{
									Type: ObjectType,
									Properties: map[string]*openapi3.SchemaRef{
										"username": {Value: &openapi3.Schema{Type: StringType}},
										"email":    {Value: &openapi3.Schema{Type: StringType}},
										"image":    {Value: &openapi3.Schema{Type: StringType}},
									},
									Required: []string{"username", "email", "image"},
								},
							},
						},
					},
				},
			},
			Responses: responses,
		},
	})

	/* user/password-change-init endpoint */
	responses = openapi3.NewResponses()
	responses.Set("201", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &descriptions[14],
		},
	})
	responses.Set("400", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &descriptions[5],
		},
	})
	responses.Set("404", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &descriptions[15],
		},
	})
	responses.Set("500", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &descriptions[7],
		},
	})
	paths.Set("/user/password-change-init", &openapi3.PathItem{
		Post: &openapi3.Operation{
			Summary: "Initialize a password change request.",
			Tags:    []string{"User"},
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Description: "Email address for password change request",
					Required:    true,
					Content: openapi3.Content{
						"application/json": &openapi3.MediaType{
							Schema: &openapi3.SchemaRef{
								Value: &openapi3.Schema{
									Type: ObjectType,
									Properties: map[string]*openapi3.SchemaRef{
										"email": {Value: &openapi3.Schema{Type: StringType}},
									},
									Required: []string{"email"},
								},
							},
						},
					},
				},
			},
			Responses: responses,
		},
	})

	/* user/change-password endpoint */
	responses = openapi3.NewResponses()
	responses.Set("200", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &descriptions[16],
		},
	})
	responses.Set("400", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &descriptions[5],
		},
	})
	responses.Set("404", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &descriptions[17],
		},
	})
	responses.Set("500", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &descriptions[7],
		},
	})
	paths.Set("/user/change-password", &openapi3.PathItem{
		Patch: &openapi3.Operation{
			Summary: "Change user's password.",
			Tags:    []string{"User"},
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Description: "Password change details",
					Required:    true,
					Content: openapi3.Content{
						"application/json": &openapi3.MediaType{
							Schema: &openapi3.SchemaRef{
								Value: &openapi3.Schema{
									Type: ObjectType,
									Properties: map[string]*openapi3.SchemaRef{
										"password": {Value: &openapi3.Schema{Type: StringType}},
										"token":    {Value: &openapi3.Schema{Type: StringType}},
									},
									Required: []string{"password", "token"},
								},
							},
						},
					},
				},
			},
			Responses: responses,
		},
	})

	/* user/password-change-valid endpoint */
	responses = openapi3.NewResponses()
	responses.Set("200", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &descriptions[18],
		},
	})
	responses.Set("410", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &descriptions[19],
		},
	})
	responses.Set("404", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &descriptions[20],
		},
	})
	responses.Set("500", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &descriptions[7],
		},
	})
	paths.Set("/user/password-change-valid", &openapi3.PathItem{
		Get: &openapi3.Operation{
			Summary: "Validate password change token.",
			Tags:    []string{"User"},
			Parameters: []*openapi3.ParameterRef{
				{
					Value: &openapi3.Parameter{
						Name:        "token",
						In:          "query",
						Description: "Password change token to validate.",
						Required:    true,
						Schema: &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type: StringType,
							},
						},
					},
				},
			},
			Responses: responses,
		},
	})

	doc := &openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:   "Lumina back-end API",
			Version: "1.0.0",
		},
		Paths: paths,
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
