package theme

import (
	"image/color"
	"math/rand"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// FlashyTheme is a custom theme with random colors
type FlashyTheme struct {
	primaryColor    color.Color
	secondaryColor  color.Color
	backgroundColor color.Color
	buttonColor     color.Color
}

// CreateRandomTheme creates a theme with random colors
func CreateRandomTheme(r *rand.Rand) fyne.Theme {
	// Generate random colors
	primaryColor := randomColor(r)
	secondaryColor := randomColor(r)
	backgroundColor := randomColor(r)
	buttonColor := randomColor(r)

	return &FlashyTheme{
		primaryColor:    primaryColor,
		secondaryColor:  secondaryColor,
		backgroundColor: backgroundColor,
		buttonColor:     buttonColor,
	}
}

// randomColor generates a random color
func randomColor(r *rand.Rand) color.Color {
	return color.RGBA{
		R: uint8(r.Intn(256)),
		G: uint8(r.Intn(256)),
		B: uint8(r.Intn(256)),
		A: 255,
	}
}

// Color implements fyne.Theme
func (t *FlashyTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNamePrimary:
		return t.primaryColor
	case "secondary":
		return t.secondaryColor
	case "background":
		return t.backgroundColor
	case "button":
		return t.buttonColor
	default:
		return theme.DefaultTheme().Color(name, variant)
	}
}

// Font implements fyne.Theme
func (t *FlashyTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

// Icon implements fyne.Theme
func (t *FlashyTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// Size implements fyne.Theme
func (t *FlashyTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

// Spacing implements fyne.Theme
func (t *FlashyTheme) Spacing() float32 {
	return 12
}

// Padding implements fyne.Theme
func (t *FlashyTheme) Padding() float32 {
	return 20
}
