package engine

import (
	"encoding/json"
	"github.com/oneaudit/go-wpjson/pkg/types"
	"github.com/oneaudit/go-wpjson/pkg/utils"
	"github.com/projectdiscovery/gologger"
	errorutil "github.com/projectdiscovery/utils/errors"
	"net/url"
	"strings"
)

type Specification struct {
	Name           string           `json:"name"`
	Description    string           `json:"description"`
	URL            string           `json:"url"`
	Home           string           `json:"home"`
	GmtOffset      int              `json:"gmt_offset"`
	TimezoneString string           `json:"timezone_string"`
	Namespaces     []string         `json:"namespaces"`
	Authentication []interface{}    `json:"authentication"`
	Routes         map[string]Route `json:"routes"`
}

type Route struct {
	Namespace string     `json:"namespace"`
	Methods   []string   `json:"methods"`
	Endpoints []Endpoint `json:"endpoints"`
}

type Endpoint struct {
	Methods    []string        `json:"methods"`
	Parameters Parameters      `json:"parameters"`
	Args       json.RawMessage `json:"args"`
}

type Parameters = map[string]Parameter

type Parameter struct {
	// Type among ["", "string", "array", "object", "boolean", "integer", ]
	Type string `json:"type"`
	// Format can be ["uri", ]
	Format   string `json:"format"`
	Pattern  string `json:"pattern"`
	Default  any    `json:"default"`
	Required bool   `json:"required"`

	// Metadata may be empty
	Title       string `json:"title"`
	Description string `json:"description"`

	// Integer constraints
	Minimum int `json:"minimum"`
	Maximum int `json:"maximum"`

	// String constraints
	MinLength int `json:"minLength"`
	MaxLength int `json:"maxLength"`

	// List of accepted values
	Enum  []string    `json:"enum"`
	OneOf []Parameter `json:"oneOf"`

	// It may be an array or an object
	//Properties []Parameter `json:"properties"`

	// todo: items
	// Maximum length of Items
	MaxItems int `json:"maxItems"`
}

type URLRequest struct {
	URL     string            `json:"url"`
	Methods string            `json:"method"`
	Body    string            `json:"body"`
	Headers map[string]string `json:"headers"`
}

func ParseSpecification(options *types.Options) (*Specification, error) {
	// Check if the input target is a URL or a file
	_, err := url.ParseRequestURI(options.InputTarget)
	var content []byte
	if err != nil {
		content, err = utils.ReadFile(options.InputTarget)
	} else {
		content, err = utils.ReadFromURL(options.InputTarget)
	}
	if err != nil {
		return nil, err
	}

	var api Specification
	err = json.Unmarshal(content, &api)
	if err != nil {
		return nil, errorutil.NewWithErr(err).Msgf("could not parse json file")
	}
	return &api, nil
}

func ParseEndpoints(options *types.Options) (endpoints []URLRequest, error error) {
	api, err := ParseSpecification(options)
	if err != nil {
		return
	}

	for path, routes := range api.Routes {
		for _, endpoint := range routes.Endpoints {
			var parameters Parameters
			err := json.Unmarshal(endpoint.Args, &parameters)
			if err != nil {
				if strings.HasPrefix(err.Error(), "json: cannot unmarshal array") {
					endpoint.Parameters = make(Parameters)
				} else {
					return nil, errorutil.NewWithErr(err).Msgf("could not parse json route %s", path)
				}
			} else {
				endpoint.Parameters = parameters
			}

			for _, method := range endpoint.Methods {
				request := URLRequest{
					URL:     "/wp-json" + utils.ExtractURLPathParameters(path),
					Methods: method,
					Body:    "",
					Headers: make(map[string]string),
				}

				query := strings.Builder{}
				rawBody := make(map[string]interface{})
				for parameterName, parameter := range parameters {
					var value any
					if parameter.Default != nil {
						value = parameter.Default
					} else if parameter.Enum != nil && len(parameter.Enum) > 0 {
						value = parameter.Enum[0]
					} else {
						value = nil
						switch parameter.Type {
						case "string":
							value = "string"
						case "number":
							value = parameter.Minimum
						case "integer":
							value = parameter.Minimum
						case "array":
							value = []interface{}{}
						case "object":
							value = map[string]interface{}{}
						case "boolean":
							value = false
						case "":
							value = ""
						default:
							gologger.Warning().Msgf("Found unexpected type %s for %s[%s]", parameter.Type, path, parameterName)
							value = ""
						}
					}

					if method == "GET" {
						payload, _ := json.Marshal(value)
						query.WriteRune('&')
						query.WriteString(parameterName)
						query.WriteRune('=')
						query.WriteString(strings.ReplaceAll(string(payload), "\"", ""))
					} else {
						rawBody[parameterName] = value
					}
				}
				if query.Len() > 0 {
					queryString := query.String()
					queryString = "?" + queryString[1:]
					request.URL = request.URL + queryString
				}
				payload, err := json.Marshal(rawBody)
				if err == nil {
					request.Headers["Content-Type"] = "application/json"
					request.Body = string(payload)
					if request.Body == "{}" {
						request.Body = ""
					}
				}

				endpoints = append(endpoints, request)
			}
		}
	}

	return
}
