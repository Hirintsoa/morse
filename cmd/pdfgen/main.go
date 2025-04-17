package main

import (
	"deliveries-pdf/internal/pdf"
	"deliveries-pdf/internal/theme"
	"fmt"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// Create a local random generator
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	myApp := app.New()
	myWindow := myApp.NewWindow("Kojakojan'ny livreur")
	myWindow.Resize(fyne.NewSize(1000, 700))

	myApp.Settings().SetTheme(theme.CreateRandomTheme(r))

	zoneEntry := widget.NewEntry()
	zoneEntry.SetPlaceHolder("Mankaiza mankaiza zoky ?")
	zoneEntry.TextStyle = fyne.TextStyle{Bold: true}

	contentEntry := widget.NewMultiLineEntry()
	contentEntry.SetPlaceHolder("Merci monsieur la Parole de m'avoir donné le Jury")
	contentEntry.Wrapping = fyne.TextWrapWord

	formContainer := container.NewVBox(
		widget.NewLabelWithStyle("Trasy:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		zoneEntry,
		widget.NewLabelWithStyle("Colleo eto le tany @ rossy:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	)

	contentContainer := container.NewVBox(contentEntry)
	contentContainer.Resize(fyne.NewSize(0, 890))

	// Create a container for the button
	buttonContainer := container.NewHBox(
		layout.NewSpacer(),
		widget.NewButton("Avoay fa maika e!", func() {
			zone := zoneEntry.Text
			content := contentEntry.Text

			if zone == "" || content == "" {
				dialog.ShowError(fmt.Errorf("mba fenoy tsara pr aloha (par respect)"), myWindow)
				return
			}

			entries := pdf.ParseContent(content)
			err := pdf.GeneratePDF(zone, entries, pdf.DefaultConfig())
			if err != nil {
				dialog.ShowError(err, myWindow)
				return
			}

			dialog.ShowInformation("Poinsa", "Tadiavo rery ao amzay", myWindow)
		}),
		layout.NewSpacer(),
	)

	// Combine all containers
	mainContainer := container.NewVBox(
		formContainer,
		contentContainer,
		buttonContainer,
	)

	// Add padding around the form
	paddedContainer := container.NewPadded(mainContainer)

	myWindow.SetContent(paddedContainer)
	myWindow.ShowAndRun()
}
