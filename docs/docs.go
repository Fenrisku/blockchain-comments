// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "https://github.com",
        "contact": {
            "name": "Fenrisku",
            "url": "https://github.com",
            "email": "fenrisku@163.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/classcom": {
            "get": {
                "description": "/classcom?cid=\u0026start=\u0026size=",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "获取评价信息"
                ],
                "summary": "获取课程评价",
                "parameters": [
                    {
                        "type": "string",
                        "description": "cid",
                        "name": "cid",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "开始数",
                        "name": "start",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "查询数量",
                        "name": "size",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"msg\": \"SUCCESS\"}",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "{\"msg\": \"FAIL}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/classinfo": {
            "get": {
                "description": "/tracecom?cid=\u0026start=\u0026size=",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "获取课程信息"
                ],
                "summary": "获取课程信息",
                "parameters": [
                    {
                        "type": "string",
                        "description": "开始数",
                        "name": "start",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "查询数量",
                        "name": "size",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"msg\": \"SUCCESS\"}",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "{\"msg\": \"FAIL}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/count": {
            "get": {
                "description": "/count",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "获取评分"
                ],
                "summary": "获取每门课程评分",
                "responses": {
                    "200": {
                        "description": "{\"msg\": \"SUCCESS\"}",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "{\"msg\": \"FAIL}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/count/total": {
            "get": {
                "description": "/count/toal",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "获取评分"
                ],
                "summary": "获取总评分",
                "responses": {
                    "200": {
                        "description": "{\"msg\": \"SUCCESS\"}",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "{\"msg\": \"FAIL}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/tracecom": {
            "get": {
                "description": "tracecom?start=\u0026size=",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "获取评价信息"
                ],
                "summary": "评价数据不分类查询",
                "parameters": [
                    {
                        "type": "string",
                        "description": "开始数",
                        "name": "start",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "查询数量",
                        "name": "size",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"msg\": \"SUCCESS\"}",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "{\"msg\": \"FAIL}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "chen-v1:8000",
	BasePath:    "/",
	Schemes:     []string{},
	Title:       "Comments-Blockchain API",
	Description: "This is a server for comments data of blockchain system.",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
		"escape": func(v interface{}) string {
			// escape tabs
			str := strings.Replace(v.(string), "\t", "\\t", -1)
			// replace " with \", and if that results in \\", replace that with \\\"
			str = strings.Replace(str, "\"", "\\\"", -1)
			return strings.Replace(str, "\\\\\"", "\\\\\\\"", -1)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
