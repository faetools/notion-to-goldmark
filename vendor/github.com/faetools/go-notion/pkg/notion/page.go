package notion

import "fmt"

// Title returns the page title.
func (p Page) Title() string {
	return p.Properties.title()
}

// Title returns the title of the page.
func (props PropertyValueMap) title() string {
	for _, prop := range props {
		if prop.Title != nil {
			return prop.Title.Content()
		}
	}

	return ""
}

// TitleWithEmoji returns the page title, prepended by an emoji if present.
func (p Page) TitleWithEmoji() string {
	if p.Icon != nil && p.Icon.Type == IconTypeEmoji {
		return fmt.Sprintf("%s %s", *p.Icon.Emoji, p.Title())
	}

	return p.Title()
}
