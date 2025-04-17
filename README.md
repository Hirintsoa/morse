# Deliveries PDF Generator

A simple application to generate PDF delivery lists for couriers.

## Features

- Generate PDF delivery lists with customer information
- Display items with prices in a clean format
- Include notes and delivery instructions
- Save PDFs to the Downloads folder with timestamps

## Installation

```bash
# Clone the repository
git clone https://github.com/Hirintsoa/morse.git
cd morse

# Install dependencies
go mod download

# Build the application
./build.sh
```

## Usage

1. Run the application:
   ```bash
   ./pdfgen
   ```

2. Enter the zone name in the "Faritra" field.

3. Enter the delivery data in the "Content" field using the following format:
   ```
   ID	Name	Address	Phone	Items	Notes
   ```
   
   Where:
   - `ID`: Customer ID
   - `Name`: Customer name
   - `Address`: Delivery address
   - `Phone`: Customer phone number
   - `Items`: List of items with prices, separated by "+" (e.g., "18+25+30")
   - `Notes`: Optional delivery notes

4. Click "Generate PDF" to create the PDF file.

5. The PDF will be saved to your Downloads folder with a timestamp.

## PDF Format

The generated PDF includes:
- Zone header
- Customer information (name, ID, address, phone)
- Items with prices (in thousands format)
- Total amount (in full Ariary format)
- Notes section
- Delivery notes box
- Date at the bottom

## License

MIT
