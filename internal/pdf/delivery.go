package pdf

import (
	"fmt"
	"strings"
)

// DeliveryEntry represents a single delivery entry with customer information and items
type DeliveryEntry struct {
	ID      string
	Name    string
	Address string
	Phone   string
	Items   string
	Notes   string
}

// ParseContent parses a tab-separated string into a slice of DeliveryEntry
func ParseContent(content string) []DeliveryEntry {
	var entries []DeliveryEntry
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		fields := strings.Split(line, "\t")
		if len(fields) >= 5 {
			// Check if all fields are empty
			allEmpty := true
			for _, field := range fields[:5] { // Check only required fields
				if strings.TrimSpace(field) != "" {
					allEmpty = false
					break
				}
			}

			if allEmpty {
				continue // Skip this entry if all fields are empty
			}

			entry := DeliveryEntry{
				ID:      fields[0],
				Name:    fields[1],
				Address: fields[2],
				Phone:   fields[3],
				Items:   fields[4],
			}
			if len(fields) > 5 {
				entry.Notes = fields[5]
			}
			entries = append(entries, entry)
		}
	}

	return entries
}

// CalculateTotal calculates the total price of all items in the delivery entry
func (e *DeliveryEntry) CalculateTotal() float64 {
	items := strings.Split(e.Items, "+")
	total := 0.0

	for _, item := range items {
		price := strings.TrimSpace(item)
		p, err := parseFloat(price)
		if err == nil {
			total += p
		}
	}

	return total
}

// parseFloat parses a string to float64 and converts to Ariary
func parseFloat(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	return f * 1000, err // Convert to Ariary
}

// FormatNumber formats a float64 as a string without decimal places
func FormatNumber(n float64) string {
	return fmt.Sprintf("%.0f", n)
}
