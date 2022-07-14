package notion

import "fmt"

// Title returns the title of the block.
func (b Block) Title() string {
	switch b.Type {
	case BlockTypeChildDatabase:
		return b.ChildDatabase.Title
	case BlockTypeChildPage:
		return b.ChildPage.Title
	default:
		return fmt.Sprintf("<no title defined for block type %q>", b.Type)
	}
}
