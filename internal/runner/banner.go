package runner

import (
	"github.com/projectdiscovery/gologger"
)

var version = "v1.0.0"

// showBanner is used to show the banner to the user
func showBanner() {
	gologger.Print().Msgf("go-wpjson %s\n", version)
}
