package engine

import (
	"fmt"
	"github.com/oneaudit/go-wpjson/pkg/types"
	"github.com/oneaudit/go-wpjson/pkg/utils"
	"net/url"
)

type WPRestEndpoint struct {
}

func ParseEndpoints(options *types.Options) (endpoints []WPRestEndpoint, error error) {
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

	fmt.Println(string(content))
	return
}
