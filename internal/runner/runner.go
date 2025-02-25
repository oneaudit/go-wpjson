package runner

import (
	"encoding/json"
	"fmt"
	"github.com/oneaudit/go-wpjson/pkg/engine"
	"github.com/oneaudit/go-wpjson/pkg/types"
	"github.com/projectdiscovery/gologger"
	errorutil "github.com/projectdiscovery/utils/errors"
)

func Execute(options *types.Options) error {
	options.ConfigureOutput()
	showBanner()

	if options.Version {
		gologger.Info().Msgf("Current version: %s", version)
		return nil
	}

	if err := validateOptions(options); err != nil {
		return errorutil.NewWithErr(err).Msgf("could not validate options")
	}

	content, err := engine.LoadContent(options)
	if err != nil {
		return errorutil.NewWithErr(err).Msgf("Could not load API content")
	}

	spec, err := engine.ParseSpecification(content)
	if err != nil {
		return errorutil.NewWithErr(err).Msgf("Could not parse API specification")
	}

	endpoints, err := engine.ParseEndpoints(spec)
	if err != nil {
		return errorutil.NewWithErr(err).Msgf("could not parse endpoints")
	}

	for _, endpoint := range endpoints {
		obj, _ := json.Marshal(endpoint)
		fmt.Println(string(obj))
	}

	return nil
}
