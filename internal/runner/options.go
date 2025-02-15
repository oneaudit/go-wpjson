package runner

import (
	"errors"
	"github.com/oneaudit/go-wpjson/pkg/types"
)

func validateOptions(options *types.Options) error {
	if options.InputTarget == "" {
		return errors.New("input file is required")
	}
	return nil
}
