package main

import (
	"fmt"
	"net/url"

	"github.com/Design-Machines-Studio/livewires-templ/form"
)

// iconPickerField is one gallery picker. Each field validates submissions
// against exactly the sections it renders, as an application must.
type iconPickerField struct {
	Name     string
	Initial  string
	Sections []form.IconPickerSection
}

var iconPickerReferenceSections = []form.IconPickerSection{{ID: "all", Heading: "Circle icons", Icons: iconPickerReferenceIcons}}

var iconPickerFields = []iconPickerField{
	{Name: "space_icon", Initial: "users-three", Sections: iconPickerCategorySections},
	{Name: "circle_icon", Sections: iconPickerReferenceSections},
	{Name: "native_circle_icon", Sections: iconPickerReferenceSections},
	{Name: "native_space_icon", Initial: "house-line", Sections: iconPickerCategorySections},
	{Name: "scheme_inherited_icon", Initial: "leaf", Sections: iconPickerReferenceSections},
	{Name: "scheme_light_icon", Sections: iconPickerReferenceSections},
	{Name: "scheme_dark_icon", Initial: "star", Sections: iconPickerReferenceSections},
}

// iconPickerDemo is the rendered state after an optional GET submission.
type iconPickerDemo struct {
	Values   map[string]string
	Saved    []string
	Rejected []string
}

func (d iconPickerDemo) value(name string) string { return d.Values[name] }

func iconPickerAllows(sections []form.IconPickerSection, name string) bool {
	for _, section := range sections {
		for _, icon := range section.Icons {
			if icon.Name == name {
				return true
			}
		}
	}
	return false
}

func iconPickerSubmission(query url.Values) iconPickerDemo {
	demo := iconPickerDemo{Values: map[string]string{}}
	for _, field := range iconPickerFields {
		demo.Values[field.Name] = field.Initial
		if !query.Has(field.Name) {
			continue
		}
		submitted := query.Get(field.Name)
		switch {
		case submitted == "":
			demo.Values[field.Name] = ""
			demo.Saved = append(demo.Saved, fmt.Sprintf("%s: no icon.", field.Name))
		case iconPickerAllows(field.Sections, submitted):
			demo.Values[field.Name] = submitted
			demo.Saved = append(demo.Saved, fmt.Sprintf("%s: %s.", field.Name, submitted))
		default:
			demo.Rejected = append(demo.Rejected, fmt.Sprintf("%s: %q is not one of this picker’s icons. Nothing was saved.", field.Name, submitted))
		}
	}
	return demo
}
