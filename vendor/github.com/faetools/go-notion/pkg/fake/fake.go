package fake

import (
	"embed" // fake responses

	"github.com/faetools/go-notion/pkg/notion"
)

// PageID is the ID of our example page.
// It can be viewed here: https://ancient-gibbon-2cd.notion.site/Example-Page-96245c8f178444a482ad1941127c3ec3
const PageID notion.Id = "96245c8f-1784-44a4-82ad-1941127c3ec3"

// Responses contains a number of responses we have generated.
//
//go:embed v1
var responses embed.FS

// HTMLExport contains the export of the example page as HTML.
//
//go:embed html
var HTMLExport embed.FS

// MDCSVExport contains the export of the example page as markdown and CSV.
//
//go:embed md-csv
var MDCSVExport embed.FS
