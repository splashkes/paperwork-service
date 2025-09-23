# PDF Template Guide for Designers

## Overview

This system allows designers to create beautiful PDF templates that the application will fill with dynamic data. You design the layout, we handle the data!

## Directory Structure

```
templates/pdf/
├── templates/          # Your PDF template files go here
│   ├── artist-page.pdf
│   ├── artist-list.pdf
│   └── auction-info.pdf
└── configs/           # Configuration files
    └── template-config.json
```

## How It Works

1. **You create a PDF template** in your favorite design tool (InDesign, Illustrator, etc.)
2. **Leave blank spaces** where dynamic content will go (artist names, QR codes, etc.)
3. **Export as PDF** and save in the `templates/` directory
4. **Update the configuration** to specify where content should be placed

## Template Types

### 1. Artist Page Template (`artist-page.pdf`)
Individual page for each artist with:
- Artist name
- Round/Easel badge
- Instagram handle
- Location
- QR code (60mm x 60mm)
- Event history table

### 2. Artist List Template (`artist-list.pdf`)
Summary page with table of all artists:
- Table with columns: #, Name, Instagram, Round, Easel
- Can have decorative headers/footers

### 3. Auction Info Template (`auction-info.pdf`)
Bidding information page:
- Table with auction/bidding data
- Columns: Art ID, Artist Name, # Bids, Top Bid, Bidder info

## Configuration File

The `template-config.json` file tells the system where to place content on your templates.

### Example Configuration

```json
{
  "artistPage": {
    "template": "artist-page.pdf",
    "fields": {
      "artistName": {
        "x": 105,        // X position in mm from left
        "y": 120,        // Y position in mm from top
        "fontSize": 36,
        "color": {"r": 41, "g": 128, "b": 185}
      },
      "qrCode": {
        "x": 20,
        "y": 100,
        "width": 60,
        "height": 60
      }
    }
  }
}
```

## Design Guidelines

### Page Size
- Use **Letter size landscape** (279mm x 216mm)
- Keep important content within margins

### Fonts
- System uses Arial by default
- Embedded fonts in PDFs are preserved for static elements

### Colors
- Art Battle Blue: RGB(41, 128, 185)
- Accent Red: RGB(231, 76, 60)
- Success Green: RGB(46, 204, 113)
- Dark Blue: RGB(52, 73, 94)

### Content Areas
Leave these areas blank in your template:

1. **QR Code Area**: 60mm x 60mm square
2. **Artist Name**: Allow ~150mm width for long names
3. **Event History**: 240mm x 50mm area for table
4. **Tables**: Full width minus margins

## Tips

1. **Test with sample data** - Some artists have long names!
2. **Use high contrast** - PDFs may be printed in B&W
3. **Keep it clean** - Too many design elements can clash with dynamic content
4. **Version your templates** - Save as `artist-page-v2.pdf` when making changes

## Updating Templates

1. Save your new PDF in the `templates/` directory
2. Update `template-config.json` if you moved any content areas
3. Test with a few events before deploying

## Need Help?

- Check if content is appearing in the wrong place? Update X/Y coordinates in config
- Text too small/large? Adjust fontSize in config
- Want to add new dynamic content? Contact the development team

## Current Dynamic Fields

### Artist Page
- `artistName` - Full name
- `roundEaselBadge` - "ROUND X | EASEL Y"
- `instagram` - Instagram handle
- `location` - City, Province
- `qrCode` - QR code image
- `qrUrl` - URL text below QR
- `eventHistory` - Past events table

### Tables
Tables are generated programmatically but respect your template's style and positioning.