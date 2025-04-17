package pdf

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

// FontPaths contains paths to regular and bold font files
type FontPaths struct {
	Regular string
	Bold    string
}

// FindFont searches for suitable fonts in common locations
func FindFont() (*FontPaths, error) {
	// Add OS-specific font paths
	switch runtime.GOOS {
	case "windows":
		// Check for Arial
		if _, err := os.Stat(`C:\Windows\Fonts\Arial.ttf`); err == nil {
			if _, err := os.Stat(`C:\Windows\Fonts\Arialbd.ttf`); err == nil {
				return &FontPaths{
					Regular: `C:\Windows\Fonts\Arial.ttf`,
					Bold:    `C:\Windows\Fonts\Arialbd.ttf`,
				}, nil
			}
		}

		// Check for Calibri
		if _, err := os.Stat(`C:\Windows\Fonts\Calibri.ttf`); err == nil {
			if _, err := os.Stat(`C:\Windows\Fonts\Calibrib.ttf`); err == nil {
				return &FontPaths{
					Regular: `C:\Windows\Fonts\Calibri.ttf`,
					Bold:    `C:\Windows\Fonts\Calibrib.ttf`,
				}, nil
			}
		}

		// Check for Segoe UI
		if _, err := os.Stat(`C:\Windows\Fonts\Segoeui.ttf`); err == nil {
			if _, err := os.Stat(`C:\Windows\Fonts\Segoeuib.ttf`); err == nil {
				return &FontPaths{
					Regular: `C:\Windows\Fonts\Segoeui.ttf`,
					Bold:    `C:\Windows\Fonts\Segoeuib.ttf`,
				}, nil
			}
		}

	case "linux":
		// Check for DejaVu Sans
		if _, err := os.Stat("/usr/share/fonts/TTF/DejaVuSans.ttf"); err == nil {
			if _, err := os.Stat("/usr/share/fonts/TTF/DejaVuSans-Bold.ttf"); err == nil {
				return &FontPaths{
					Regular: "/usr/share/fonts/TTF/DejaVuSans.ttf",
					Bold:    "/usr/share/fonts/TTF/DejaVuSans-Bold.ttf",
				}, nil
			}
		}

		// Check for Liberation Sans
		if _, err := os.Stat("/usr/share/fonts/truetype/liberation/LiberationSans-Regular.ttf"); err == nil {
			if _, err := os.Stat("/usr/share/fonts/truetype/liberation/LiberationSans-Bold.ttf"); err == nil {
				return &FontPaths{
					Regular: "/usr/share/fonts/truetype/liberation/LiberationSans-Regular.ttf",
					Bold:    "/usr/share/fonts/truetype/liberation/LiberationSans-Bold.ttf",
				}, nil
			}
		}

		// Check for Noto Sans
		if _, err := os.Stat("/usr/share/fonts/noto/NotoSans-Regular.ttf"); err == nil {
			if _, err := os.Stat("/usr/share/fonts/noto/NotoSans-Bold.ttf"); err == nil {
				return &FontPaths{
					Regular: "/usr/share/fonts/noto/NotoSans-Regular.ttf",
					Bold:    "/usr/share/fonts/noto/NotoSans-Bold.ttf",
				}, nil
			}
		}

	case "darwin":
		// Check for Arial
		if _, err := os.Stat("/Library/Fonts/Arial.ttf"); err == nil {
			if _, err := os.Stat("/Library/Fonts/Arial Bold.ttf"); err == nil {
				return &FontPaths{
					Regular: "/Library/Fonts/Arial.ttf",
					Bold:    "/Library/Fonts/Arial Bold.ttf",
				}, nil
			}
		}

		// Check for Helvetica
		if _, err := os.Stat("/System/Library/Fonts/Helvetica.ttc"); err == nil {
			// Helvetica.ttc is a collection, we'll use it for both regular and bold
			return &FontPaths{
				Regular: "/System/Library/Fonts/Helvetica.ttc",
				Bold:    "/System/Library/Fonts/Helvetica.ttc",
			}, nil
		}
	}

	// If no system fonts found, try to use fallback fonts
	return SetupFallbackFonts()
}

// SetupFallbackFonts downloads Liberation Sans fonts if needed
func SetupFallbackFonts() (*FontPaths, error) {
	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not get home directory: %v", err)
	}

	// Create fonts directory in user's home directory
	fontsDir := filepath.Join(homeDir, ".fonts")
	if err := os.MkdirAll(fontsDir, 0755); err != nil {
		return nil, fmt.Errorf("could not create fonts directory: %v", err)
	}

	// Use Liberation Sans as embedded fallback
	regularFontPath := filepath.Join(fontsDir, "LiberationSans-Regular.ttf")
	boldFontPath := filepath.Join(fontsDir, "LiberationSans-Bold.ttf")

	if _, err := os.Stat(regularFontPath); err != nil {
		// Download Liberation Sans if not present
		if err := downloadLiberationSans(regularFontPath); err != nil {
			return nil, fmt.Errorf("could not setup fallback font: %v", err)
		}
	}

	if _, err := os.Stat(boldFontPath); err != nil {
		// Download Liberation Sans Bold if not present
		if err := downloadLiberationSansBold(boldFontPath); err != nil {
			return nil, fmt.Errorf("could not setup fallback bold font: %v", err)
		}
	}

	return &FontPaths{
		Regular: regularFontPath,
		Bold:    boldFontPath,
	}, nil
}

// downloadLiberationSans downloads the Liberation Sans Regular font
func downloadLiberationSans(destPath string) error {
	// URL for Liberation Sans Regular
	fontURL := "https://github.com/liberationfonts/liberation-fonts/raw/main/liberation-fonts-ttf-2.1.5/LiberationSans-Regular.ttf"

	// Download the font file
	resp, err := http.Get(fontURL)
	if err != nil {
		return fmt.Errorf("failed to download font: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download font: HTTP status %d", resp.StatusCode)
	}

	// Create the destination file
	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create font file: %v", err)
	}
	defer out.Close()

	// Copy the content
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save font file: %v", err)
	}

	return nil
}

// downloadLiberationSansBold downloads the Liberation Sans Bold font
func downloadLiberationSansBold(destPath string) error {
	// URL for Liberation Sans Bold
	fontURL := "https://github.com/liberationfonts/liberation-fonts/raw/main/liberation-fonts-ttf-2.1.5/LiberationSans-Bold.ttf"

	// Download the font file
	resp, err := http.Get(fontURL)
	if err != nil {
		return fmt.Errorf("failed to download font: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download font: HTTP status %d", resp.StatusCode)
	}

	// Create the destination file
	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create font file: %v", err)
	}
	defer out.Close()

	// Copy the content
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save font file: %v", err)
	}

	return nil
}
