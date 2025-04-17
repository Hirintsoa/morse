#!/bin/bash

# Create assets directory if it doesn't exist
mkdir -p assets

# Check if logo exists, if not create a placeholder
if [ ! -f "assets/logo.png" ]; then
    echo "Please add your logo.png file to the assets directory"
    echo "The logo should be a PNG file, ideally 200x200 pixels"
fi

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o pdfgen ./cmd/pdfgen

# Build for Windows with icon
if [ -f "assets/logo.png" ]; then
    # Convert PNG to ICO for Windows
    magick assets/logo.png -define icon:auto-resize=16,32,48,64,128,256 assets/logo.ico
    
    # Build Windows executable with icon
    GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -ldflags="-H windowsgui" -o pdfgen.exe ./cmd/pdfgen
    
    # Add icon to Windows executable
    if command -v rsrc &> /dev/null; then
        rsrc -ico assets/logo.ico -o rsrc.syso
        GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -ldflags="-H windowsgui" -o pdfgen.exe ./cmd/pdfgen
        rm rsrc.syso
    else
        echo "rsrc not found. Please install it with: go install github.com/akavel/rsrc@latest"
    fi
else
    # Build Windows executable without icon
    GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -ldflags="-H windowsgui" -o pdfgen.exe ./cmd/pdfgen
fi

echo "Build complete!" 