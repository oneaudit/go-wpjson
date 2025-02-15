package engine

import (
	"encoding/json"
	"github.com/oneaudit/go-wpjson/pkg/types"
	"github.com/oneaudit/go-wpjson/pkg/utils"
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

func ParseEndpoints(options *types.Options) (api *Specification, error error) {
	// Check if the input target is a URL or a file
	_, error = url.ParseRequestURI(options.InputTarget)
	var content []byte
	if error != nil {
		content, error = utils.ReadFile(options.InputTarget)
	} else {
		content, error = utils.ReadFromURL(options.InputTarget)
	}
	if error != nil {
		return
	}

	err := json.Unmarshal(content, &api)
	if err != nil {
		return nil, errorutil.NewWithErr(err).Msgf("could not parse json file")
	}

	for path, routes := range api.Routes {
		for _, endpoint := range routes.Endpoints {
			var parameters Parameters
			err := json.Unmarshal(endpoint.Args, &parameters)
			if err != nil {
				if strings.HasPrefix(err.Error(), "json: cannot unmarshal array") {
					endpoint.Parameters = make(Parameters)
					continue
				}
				return nil, errorutil.NewWithErr(err).Msgf("could not parse json route %s", path)
			} else {
				endpoint.Parameters = parameters
			}
		}
	}

	return
}
