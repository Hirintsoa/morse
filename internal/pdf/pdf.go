package pdf

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/signintech/gopdf"
)

type PDFConfig struct {
	PageWidth      float64
	MarginLeft     float64
	MarginRight    float64
	MarginTop      float64
	MarginBottom   float64
	LineHeight     float64
	SectionSpacing float64
	ItemSpacing    float64
	NoteBoxHeight  float64
	NameSpacing    float64
	EntrySpacing   float64
	ItemWidth      float64
	DateSpacing    float64
	PhoneSpacing   float64
	ZoneSpacing    float64
	AddressSpacing float64
}

func DefaultConfig() *PDFConfig {
	return &PDFConfig{
		PageWidth:      78.0,
		MarginLeft:     4.0,
		MarginRight:    4.0,
		MarginTop:      4.0,
		MarginBottom:   2.0,
		LineHeight:     2.0,
		SectionSpacing: 2.0,
		ItemSpacing:    2.5,
		NoteBoxHeight:  8.0,
		NameSpacing:    2.0,
		EntrySpacing:   4.0,
		ItemWidth:      25.0,
		DateSpacing:    3.0,
		PhoneSpacing:   2.0,
		ZoneSpacing:    3.0,
		AddressSpacing: 1.0,
	}
}

func GeneratePDF(zone string, entries []DeliveryEntry, config *PDFConfig) error {
	if config == nil {
		config = DefaultConfig()
	}

	pdf := gopdf.GoPdf{}
	totalHeight := (float64(len(entries)) * 50) + 30
	pdf.Start(gopdf.Config{
		PageSize: gopdf.Rect{W: config.PageWidth, H: totalHeight},
		Unit:     gopdf.Unit_MM,
	})

	fontPaths, err := FindFont()
	if err != nil {
		fontPaths, err = SetupFallbackFonts()
		if err != nil {
			return fmt.Errorf("could not setup fonts: %v", err)
		}
	}

	err = pdf.AddTTFFont("regular", fontPaths.Regular)
	if err != nil {
		return fmt.Errorf("could not load regular font %s: %v", fontPaths.Regular, err)
	}

	err = pdf.AddTTFFont("bold", fontPaths.Bold)
	if err != nil {
		return fmt.Errorf("could not load bold font %s: %v", fontPaths.Bold, err)
	}

	pdf.AddPage()
	currentY := config.MarginTop

	logoPath := "assets/logo.png"
	if _, err := os.Stat(logoPath); err == nil {
		pdf.Image(logoPath, (config.PageWidth-20)/2, currentY, &gopdf.Rect{W: 18, H: 16})
		currentY += 12
	}

	// Add Zone header
	pdf.SetFont("regular", "", 10)
	pdf.SetX(config.MarginLeft)
	pdf.SetY(currentY)
	pdf.SetFont("bold", "", 12)

	// Calculate available width for zone text
	availableWidth := config.PageWidth - config.MarginLeft - config.MarginRight - 4
	zoneLines := wrapText(&pdf, zone, availableWidth)
	for _, line := range zoneLines {
		pdf.SetX(config.MarginLeft + 2)
		pdf.SetY(currentY)
		pdf.Cell(nil, line)
		currentY += config.LineHeight + 1.0
	}
	currentY += config.ZoneSpacing

	for i, entry := range entries {
		pdf.SetFont("bold", "", 9)

		total := entry.CalculateTotal()

		// Header with customer name and total on the right
		pdf.SetX(config.MarginLeft)
		pdf.SetY(currentY)
		pdf.Cell(nil, entry.Name)

		// Total on the right
		totalText := fmt.Sprintf("%s Ar", FormatNumber(total))
		totalWidth, _ := pdf.MeasureTextWidth(totalText)
		pdf.SetX(config.PageWidth - config.MarginRight - totalWidth)
		pdf.Cell(nil, totalText)

		currentY += config.LineHeight + 2.0

		pdf.SetFont("regular", "", 7)
		pdf.SetX(config.MarginLeft)
		pdf.SetY(currentY)
		if entry.ID == "" {
			pdf.Cell(nil, "ID: -")
		} else {
			pdf.Cell(nil, fmt.Sprintf("ID: %s", entry.ID))
		}
		currentY += config.LineHeight + config.NameSpacing

		// Address with wrapping
		pdf.SetFont("regular", "", 8)
		addressLines := wrapText(&pdf, entry.Address, config.PageWidth-config.MarginLeft-config.MarginRight-7) // -7 for icon and spacing
		for _, line := range addressLines {
			pdf.SetX(config.MarginLeft + 6)
			pdf.SetY(currentY)
			if line == addressLines[0] {
				pdf.SetX(config.MarginLeft)
				pdf.Cell(nil, ">")
				pdf.SetX(config.MarginLeft + 4)
			}
			pdf.Cell(nil, line)
			currentY += config.LineHeight + 0.5
		}
		currentY += config.AddressSpacing

		// Phone number
		pdf.SetX(config.MarginLeft)
		pdf.SetY(currentY)
		pdf.Cell(nil, "#")
		pdf.SetX(config.MarginLeft + 4)
		if entry.Phone == "" {
			pdf.Cell(nil, "Tsisy lty a! Tsisy")
		} else {
			pdf.Cell(nil, entry.Phone)
		}
		currentY += config.LineHeight + config.PhoneSpacing

		// Items section
		pdf.SetX(config.MarginLeft)
		pdf.SetY(currentY)
		pdf.Cell(nil, "Entam-be:")
		currentY += config.LineHeight + config.ItemSpacing

		items := strings.Split(entry.Items, "+")
		for i := 0; i < len(items); i += 3 {
			for j := 0; j < 3 && i+j < len(items); j++ {
				price := strings.TrimSpace(items[i+j])
				p, err := parseFloat(price)
				if err == nil {
					pdf.SetX(config.MarginLeft + (config.ItemWidth * float64(j)))
					pdf.SetY(currentY)
					if p == 0 {
						pdf.Cell(nil, "• Kadoa")
					} else {
						pdf.Cell(nil, fmt.Sprintf("• %sk", FormatNumber(p/1000)))
					}
				}
			}
			currentY += config.LineHeight + config.ItemSpacing
		}

		// Notes from customer
		if entry.Notes != "" {
			currentY += config.SectionSpacing
			pdf.SetFont("bold", "", 8)
			pdf.SetX(config.MarginLeft)
			pdf.SetY(currentY)
			pdf.Cell(nil, "Notes:")
			pdf.SetX(config.MarginLeft + 12)
			pdf.SetFont("regular", "", 8)

			// Calculate available width for notes text
			notesWidth := config.PageWidth - config.MarginLeft - config.MarginRight - 12
			noteLines := wrapText(&pdf, entry.Notes, notesWidth)
			for _, line := range noteLines {
				pdf.SetX(config.MarginLeft + 12)
				pdf.SetY(currentY)
				pdf.Cell(nil, line)
				currentY += config.LineHeight + 0.5
			}
		}

		// Note-taking box for deliverer
		currentY += config.SectionSpacing
		pdf.SetLineWidth(0.1)
		pdf.SetLineType("dashed")
		noteBoxY := currentY
		pdf.Line(config.MarginLeft, noteBoxY, config.PageWidth-config.MarginRight, noteBoxY)                                           // Top
		pdf.Line(config.MarginLeft, noteBoxY+config.NoteBoxHeight, config.PageWidth-config.MarginRight, noteBoxY+config.NoteBoxHeight) // Bottom
		pdf.Line(config.MarginLeft, noteBoxY, config.MarginLeft, noteBoxY+config.NoteBoxHeight)                                        // Left
		pdf.Line(config.PageWidth-config.MarginRight, noteBoxY, config.PageWidth-config.MarginRight, noteBoxY+config.NoteBoxHeight)    // Right

		pdf.SetFont("regular", "", 6)
		pdf.SetFont("bold", "", 8)
		pdf.SetX(config.MarginLeft + 1)
		pdf.SetY(noteBoxY + 2)
		pdf.Cell(nil, "Watawata:")
		currentY += config.NoteBoxHeight

		pdf.SetLineType("solid")
		if i < len(entries)-1 {
			currentY += config.EntrySpacing / 2
			pdf.SetLineWidth(0.3)
			pdf.Line(config.MarginLeft, currentY, config.PageWidth-config.MarginRight, currentY)
			currentY += config.EntrySpacing / 2
		}
	}

	// Add current date
	currentY += config.DateSpacing
	pdf.SetFont("regular", "", 8)
	dateStr := time.Now().Format("02/01/2006")
	dateWidth, _ := pdf.MeasureTextWidth(dateStr)
	pdf.SetX((config.PageWidth - dateWidth) / 2)
	pdf.SetY(currentY)
	pdf.Cell(nil, dateStr)

	currentY += config.LineHeight + 2.0
	quote := "\"Taloha sarotra nirahana, ankehitriny lasa livreur.🥲\""
	quoteWidth, _ := pdf.MeasureTextWidth(quote)
	pdf.SetFont("bold", "I", 10)
	pdf.SetX((config.PageWidth - quoteWidth) / 2)
	pdf.SetY(currentY)
	pdf.Cell(nil, quote)

	// Add the final word to this non-sense
	currentY += config.LineHeight + 2.0
	slogan := "KIMBASÔ !"
	sloganWidth, _ := pdf.MeasureTextWidth(slogan)
	pdf.SetFont("bold", "", 14)
	pdf.SetFont("regular", "", 11)
	pdf.SetX((config.PageWidth - sloganWidth) / 2)
	pdf.SetY(currentY)
	pdf.Cell(nil, slogan)

	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not get home directory: %v", err)
	}

	downloadsDir := filepath.Join(homeDir, "Downloads")

	if err := os.MkdirAll(downloadsDir, 0755); err != nil {
		return fmt.Errorf("could not create Downloads directory: %v", err)
	}

	timestamp := time.Now().Format("2006-01-02")
	filename := filepath.Join(downloadsDir, fmt.Sprintf("fanatitra_%s_%s.pdf", zone, timestamp))

	return pdf.WritePdf(filename)
}

func wrapText(pdf *gopdf.GoPdf, text string, maxWidth float64) []string {
	var lines []string
	words := strings.Fields(text)
	currentLine := ""

	for _, word := range words {
		testLine := currentLine
		if testLine != "" {
			testLine += " "
		}
		testLine += word

		width, _ := pdf.MeasureTextWidth(testLine)
		if width > maxWidth {
			// Line would be too long
			if currentLine != "" {
				// If we already have content, start a new line
				lines = append(lines, currentLine)
				currentLine = word
			} else {
				// For long words, split them into smaller chunks
				remainingWord := word
				for len(remainingWord) > 0 {
					// Try to find a good breaking point
					splitPoint := len(remainingWord)
					for i := 0; i < len(remainingWord); i++ {
						testChunk := remainingWord[:i+1]
						width, _ := pdf.MeasureTextWidth(testChunk)
						if width > maxWidth {
							splitPoint = i
							break
						}
					}

					// Add the chunk to lines
					lines = append(lines, remainingWord[:splitPoint])
					remainingWord = remainingWord[splitPoint:]
				}
				currentLine = ""
			}
		} else {
			// Add word to current line
			currentLine = testLine
		}
	}

	// Add the last line if not empty
	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}
