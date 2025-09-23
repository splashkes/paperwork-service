# Art Battle Paperwork Service

A standalone Go microservice that generates professional PDF paperwork for Art Battle events. This service integrates with Supabase to fetch event and artist data, then generates comprehensive PDF documents for event management.

## Features

- **Supabase Integration**: Fetches event and artist data via Supabase Edge Functions
- **Professional PDF Generation**: Creates multi-page PDF documents with custom fonts and backgrounds
- **QR Code Generation**: Dynamic QR codes linking to artist Instagram profiles or event pages
- **Event History**: Displays artist participation history with winner indicators
- **Clean Layout**: No section headings, larger readable fonts (14pt bio, 12pt history)
- **RESTful API**: Simple HTTP endpoints for PDF generation

## Architecture

- **Data Source**: Supabase database via Edge Functions
- **PDF Engine**: jung-kurt/gofpdf with custom Acumin Pro fonts
- **Background Images**: Designer-provided templates for professional appearance
- **QR Codes**: Dynamic generation with Instagram priority, event page fallback

## API Endpoints

- `GET /api/v1/health` - Health check
- `GET /api/v1/event-pdf/{eid}` - Generate PDF for event (e.g., AB2940)

## Quick Start

```bash
# Build the service
go build -o bin/paperwork-service cmd/main.go

# Run locally
./bin/paperwork-service

# Generate PDF
curl -o event.pdf "http://localhost:8080/api/v1/event-pdf/AB2940"
```

## Environment Variables

```bash
SUPABASE_URL=https://your-project.supabase.co
PORT=8080
ENVIRONMENT=development
TEMPLATES_PATH=./assets
```

## Deployment

Designed for deployment to DigitalOcean App Platform or similar container platforms.

## Data Flow

1. HTTP request with event EID
2. Fetch data from Supabase Edge Function
3. Generate PDF with artist pages and event history
4. Return PDF as downloadable attachment

Built for Art Battle event management system.
