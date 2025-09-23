package services

import (
	"bytes"
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"paperwork-service/internal/models"

	"github.com/jung-kurt/gofpdf"
	"github.com/skip2/go-qrcode"
	"go.uber.org/zap"
)

// PaperworkPDFService generates PDFs with designer-provided background images
type PaperworkPDFService struct {
	logger          *zap.Logger
	backgroundsPath string
	fontsPath       string
}

// NewPaperworkPDFService creates a new background-based PDF service
func NewPaperworkPDFService(logger *zap.Logger, templatesPath string) *PaperworkPDFService {
	return &PaperworkPDFService{
		logger:          logger,
		backgroundsPath: filepath.Join(templatesPath, "backgrounds"),
		fontsPath:       filepath.Join(templatesPath, "fonts"),
	}
}

// GenerateEventPaperwork generates the PDF with background images
func (s *PaperworkPDFService) GenerateEventPaperwork(event *models.Event, artists []models.EventArtist, auctionLots []models.AuctionLot) ([]byte, error) {
	// Create PDF in landscape mode
	pdf := gofpdf.New("L", "mm", "Letter", "")
	pdf.SetAutoPageBreak(false, 0)

	// Add custom fonts
	s.addCustomFonts(pdf)

	// Add artist list page with background
	s.addPageWithBackground(pdf, "artist-list-bg.png")
	s.addArtistListContent(pdf, event.Name, artists)

	// Add auction info page with background
	s.addPageWithBackground(pdf, "auction-info-bg.png")
	s.addAuctionInfoContent(pdf, event.Name, event.EID, event.Currency, artists, auctionLots)

	// Add bio summary pages
	s.addBioSummaryPages(pdf, event.Name, artists)

	// Add individual artist pages with background (only for ready artists)
	for _, artist := range artists {
		// Skip confirmed-only artists for individual pages
		if artist.Status == "confirmed-only" {
			continue
		}
		s.addPageWithBackground(pdf, "artist-page-bg.png")
		s.addArtistPageContent(pdf, event.Name, event.EID, artist)
	}

	// Generate PDF
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return buf.Bytes(), nil
}

// addCustomFonts registers the custom TTF fonts
func (s *PaperworkPDFService) addCustomFonts(pdf *gofpdf.Fpdf) {
	// Add Acumin Pro fonts
	fonts := map[string]string{
		"AcuminMedium":   "Acumin Pro SemiCond Medium.ttf",
		"AcuminBold":     "Acumin Pro Cond Bold.ttf",
		"AcuminSemibold": "Acumin Pro SemiCond Semibold.ttf",
	}

	for fontName, fileName := range fonts {
		fontPath := filepath.Join(s.fontsPath, fileName)
		if _, err := os.Stat(fontPath); err == nil {
			pdf.AddUTF8Font(fontName, "", fontPath)
			s.logger.Debug("Added custom font",
				zap.String("name", fontName),
				zap.String("file", fileName))
		} else {
			s.logger.Warn("Font file not found, using default",
				zap.String("name", fontName),
				zap.String("file", fileName),
				zap.Error(err))
		}
	}
}

// addPageWithBackground adds a new page with a background image
func (s *PaperworkPDFService) addPageWithBackground(pdf *gofpdf.Fpdf, backgroundFile string) {
	pdf.AddPage()

	// Try to load background image
	bgPath := filepath.Join(s.backgroundsPath, backgroundFile)
	if _, err := os.Stat(bgPath); err == nil {
		// Register the image if not already registered
		pdf.RegisterImageOptions(bgPath, gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true})

		// Place background image covering full page
		pdf.ImageOptions(bgPath, 0, 0, 279.4, 215.9, false, gofpdf.ImageOptions{ImageType: "PNG"}, 0, "")

		s.logger.Debug("Added background image", zap.String("background", backgroundFile))
	} else {
		s.logger.Debug("Background image not found, using blank page",
			zap.String("background", backgroundFile),
			zap.Error(err))
	}
}

// addArtistListContent adds the artist list content
func (s *PaperworkPDFService) addArtistListContent(pdf *gofpdf.Fpdf, eventName string, artists []models.EventArtist) {
	// Add content on top of background
	pdf.SetFont("AcuminBold", "", 24)
	pdf.SetXY(20, 20)
	pdf.Cell(0, 10, cleanString(eventName))

	// Artist table starting at specific position
	pdf.SetFont("AcuminMedium", "", 10)
	pdf.SetXY(20, 40)

	// Table headers
	headers := []string{"Round-Easel", "Artist Name"}
	colWidths := []float64{40, 130}

	pdf.SetFont("AcuminSemibold", "", 10)
	pdf.SetTextColor(0, 0, 0) // Ensure black text
	pdf.SetDrawColor(200, 200, 200) // Light gray for grid

	// Draw header row with full borders
	x := pdf.GetX()
	for i, header := range headers {
		pdf.SetX(x) // Reset X position to ensure alignment
		pdf.CellFormat(colWidths[i], 8, header, "1", 0, "C", false, 0, "") // Full border
		x += colWidths[i]
	}
	pdf.Ln(8)

	// Table rows with full grid - only show ready artists
	pdf.SetFont("AcuminMedium", "", 10)
	for _, artist := range artists {
		// Skip confirmed-only artists
		if artist.Status == "confirmed-only" {
			continue
		}

		x = 20 // Reset to starting position
		pdf.SetX(x)

		roundEaselText := fmt.Sprintf("%d-%d", artist.RoundNumber, artist.EaselNumber)
		artistName := artist.DisplayName
		if artistName == "" {
			artistName = artist.ArtistName
		}
		if artistName == "" {
			artistName = fmt.Sprintf("%s %s", artist.FirstName, artist.LastName)
		}

		// Draw cells with borders
		pdf.CellFormat(colWidths[0], 8, roundEaselText, "1", 0, "C", false, 0, "")
		pdf.CellFormat(colWidths[1], 8, cleanString(artistName), "1", 0, "L", false, 0, "")
		pdf.Ln(8)
	}
}

// addAuctionInfoContent adds the auction information content
func (s *PaperworkPDFService) addAuctionInfoContent(pdf *gofpdf.Fpdf, eventName string, eventEID string, currency string, artists []models.EventArtist, auctionLots []models.AuctionLot) {
	// Add content on top of background - match original exactly
	pdf.SetFont("AcuminBold", "", 20)
	pdf.SetXY(20, 20)
	pdf.Cell(0, 10, "Auction & Bidding Information")

	// Auction table starting at specific position
	pdf.SetFont("AcuminMedium", "", 9)
	pdf.SetXY(20, 40)

	// Table headers
	headers := []string{"EID-Round-Easel", "Artist Name", "# Bids", "Top Bid", "Bidder Info", "Payment Status"}
	colWidths := []float64{40, 60, 20, 25, 60, 35}

	pdf.SetFont("AcuminSemibold", "", 9)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetDrawColor(200, 200, 200)

	// Draw header row
	x := pdf.GetX()
	for i, header := range headers {
		pdf.SetX(x)
		pdf.CellFormat(colWidths[i], 8, header, "1", 0, "C", false, 0, "")
		x += colWidths[i]
	}
	pdf.Ln(8)

	// Create a map of auction lots by round and easel for quick lookup
	lotMap := make(map[string]models.AuctionLot)
	for _, lot := range auctionLots {
		key := fmt.Sprintf("%d-%d", lot.Round, lot.EaselNumber)
		lotMap[key] = lot
	}

	// Table rows
	pdf.SetFont("AcuminMedium", "", 9)
	for _, artist := range artists {
		if artist.Status == "confirmed-only" {
			continue
		}

		x = 20
		pdf.SetX(x)

		eidRoundEasel := fmt.Sprintf("%s-%d-%d", eventEID, artist.RoundNumber, artist.EaselNumber)
		artistName := artist.DisplayName
		if artistName == "" {
			artistName = artist.ArtistName
		}
		if artistName == "" {
			artistName = fmt.Sprintf("%s %s", artist.FirstName, artist.LastName)
		}

		// Look up auction data
		lotKey := fmt.Sprintf("%d-%d", artist.RoundNumber, artist.EaselNumber)
		lot, hasLot := lotMap[lotKey]

		bidCount := "0"
		topBid := "-"
		bidderInfo := "-"
		paymentStatus := "-"

		if hasLot {
			bidCount = fmt.Sprintf("%d", lot.BidCount)
			if lot.HighestBid > 0 {
				currencySymbol := "$"
				if currency == "EUR" {
					currencySymbol = "€"
				} else if currency == "GBP" {
					currencySymbol = "£"
				}
				topBid = fmt.Sprintf("%s%.0f", currencySymbol, lot.HighestBid)
			}
			if lot.WinningBid != nil {
				bidderInfo = lot.WinningBid.BidderName
				if bidderInfo == "" {
					bidderInfo = lot.WinningBid.BidderEmail
				}
				paymentStatus = lot.WinningBid.PaymentStatus
			}
		}

		// Draw cells
		pdf.CellFormat(colWidths[0], 8, eidRoundEasel, "1", 0, "C", false, 0, "")
		pdf.CellFormat(colWidths[1], 8, cleanString(artistName), "1", 0, "L", false, 0, "")
		pdf.CellFormat(colWidths[2], 8, bidCount, "1", 0, "C", false, 0, "")
		pdf.CellFormat(colWidths[3], 8, topBid, "1", 0, "C", false, 0, "")
		pdf.CellFormat(colWidths[4], 8, cleanString(bidderInfo), "1", 0, "L", false, 0, "")
		pdf.CellFormat(colWidths[5], 8, cleanString(paymentStatus), "1", 0, "C", false, 0, "")
		pdf.Ln(8)
	}
}

// addBioSummaryPages adds bio summary pages grouped by round
func (s *PaperworkPDFService) addBioSummaryPages(pdf *gofpdf.Fpdf, eventName string, artists []models.EventArtist) {
	// Group artists by round
	round1Artists := []models.EventArtist{}
	round2Artists := []models.EventArtist{}
	unmatchedArtists := []models.EventArtist{}

	for _, artist := range artists {
		if artist.Status == "confirmed-only" {
			unmatchedArtists = append(unmatchedArtists, artist)
		} else if artist.RoundNumber == 1 {
			round1Artists = append(round1Artists, artist)
		} else if artist.RoundNumber == 2 {
			round2Artists = append(round2Artists, artist)
		}
	}

	// Add Round 1 bios if we have artists
	if len(round1Artists) > 0 {
		s.addPageWithBackground(pdf, "artist-list-bg.png")
		s.addRoundBiosContent(pdf, eventName, "Round 1 Artist Bios", round1Artists)
	}

	// Add Round 2 bios if we have artists
	if len(round2Artists) > 0 {
		s.addPageWithBackground(pdf, "artist-list-bg.png")
		s.addRoundBiosContent(pdf, eventName, "Round 2 Artist Bios", round2Artists)
	}

	// Add unmatched/confirmed-only artists if we have them
	if len(unmatchedArtists) > 0 {
		s.addPageWithBackground(pdf, "artist-list-bg.png")
		s.addRoundBiosContent(pdf, eventName, "Additional Artist Bios", unmatchedArtists)
	}
}

// addRoundBiosContent adds bio content for a specific round
func (s *PaperworkPDFService) addRoundBiosContent(pdf *gofpdf.Fpdf, eventName string, roundTitle string, artists []models.EventArtist) {
	pdf.SetFont("AcuminBold", "", 24)
	pdf.SetXY(20, 20)
	pdf.Cell(0, 10, fmt.Sprintf("%s - %s", cleanString(eventName), roundTitle))

	pdf.SetFont("AcuminMedium", "", 10)
	pdf.SetXY(20, 40)

	for _, artist := range artists {
		artistName := artist.DisplayName
		if artistName == "" {
			artistName = artist.ArtistName
		}
		if artistName == "" {
			artistName = fmt.Sprintf("%s %s", artist.FirstName, artist.LastName)
		}

		// Add artist name
		pdf.SetFont("AcuminSemibold", "", 12)
		pdf.Cell(0, 8, cleanString(artistName))
		pdf.Ln(8)

		// Add bio if available
		if artist.Bio != "" {
			pdf.SetFont("AcuminMedium", "", 10)
			bio := cleanString(artist.Bio)
			// Word wrap bio text
			pdf.MultiCell(0, 6, bio, "", "L", false)
			pdf.Ln(4)
		} else {
			pdf.SetFont("AcuminMedium", "", 10)
			pdf.Cell(0, 6, "No bio available")
			pdf.Ln(10)
		}
	}
}

// addArtistPageContent adds individual artist page content
func (s *PaperworkPDFService) addArtistPageContent(pdf *gofpdf.Fpdf, eventName string, eventEID string, artist models.EventArtist) {
	// Constants for layout - match original exactly
	const (
		pageWidth     = 279.4
		pageHeight    = 215.9
		qrX           = 20     // QR code on left side
		qrY           = 45     // QR code Y position
		qrSize        = 42     // QR code size
		nameStartX    = 70     // Artist name starts at X=70
		nameStartY    = 48     // Artist name Y position
	)

	artistName := artist.DisplayName
	if artistName == "" {
		artistName = artist.ArtistName
	}
	if artistName == "" {
		artistName = fmt.Sprintf("%s %s", artist.FirstName, artist.LastName)
	}

	// Generate QR code - prioritize Instagram, fallback to event page
	var qrURL string
	if artist.Instagram != "" {
		// Use Instagram URL if available
		if strings.HasPrefix(artist.Instagram, "http") {
			qrURL = artist.Instagram
		} else {
			qrURL = fmt.Sprintf("https://instagram.com/%s", strings.TrimPrefix(artist.Instagram, "@"))
		}
	} else {
		// Fallback to event page
		qrURL = fmt.Sprintf("https://artb.art/event/%s", eventEID)
	}

	qr, err := qrcode.New(qrURL, qrcode.Medium)
	if err == nil {
		qr.DisableBorder = true
		img := qr.Image(256)

		// Convert to PNG using original method
		var buf bytes.Buffer
		if err := png.Encode(&buf, img); err == nil {
			imgReader := bytes.NewReader(buf.Bytes())
			imgName := fmt.Sprintf("qr_%d", artist.EntryID)
			pdf.RegisterImageOptionsReader(imgName, gofpdf.ImageOptions{ImageType: "png"}, imgReader)

			// Place QR code at exact original position
			pdf.ImageOptions(imgName, qrX, qrY, qrSize, qrSize, false, gofpdf.ImageOptions{ImageType: "png"}, 0, "")
		}
	}

	// Artist name with dynamic font sizing - match original exactly
	rightMargin := float64(20)
	availableWidth := pageWidth - nameStartX - rightMargin

	// Start with default font size and check if name fits
	fontSize := float64(49) // Default size
	pdf.SetFont("AcuminBold", "", fontSize)
	nameWidth := pdf.GetStringWidth(artistName)

	// Reduce font size if name is too long
	for nameWidth > availableWidth && fontSize > 20 {
		fontSize -= 2
		pdf.SetFont("AcuminBold", "", fontSize)
		nameWidth = pdf.GetStringWidth(artistName)
	}

	pdf.SetTextColor(0, 0, 0) // Black text
	pdf.SetXY(nameStartX, nameStartY)
	pdf.Cell(availableWidth, 10, cleanString(artistName))

	// Event name above round/easel - match original positioning
	pdf.SetFont("AcuminMedium", "", 18)
	pdf.SetXY(nameStartX, 69)
	pdf.Cell(0, 6, cleanString(eventName))

	pdf.SetFont("AcuminMedium", "", 21)
	pdf.SetXY(nameStartX, 77)
	pdf.Cell(0, 8, fmt.Sprintf("Round %d - Easel %d", artist.RoundNumber, artist.EaselNumber))

	// Two-column layout for bottom half of page - match original exactly
	startY := float64(118) // Start position for both columns
	columnWidth := float64(115) // Width of each column
	columnGap := float64(10) // Gap between columns
	leftColumnX := float64(20)
	rightColumnX := leftColumnX + columnWidth + columnGap
	rightColumnMaxWidth := float64(279.4 - rightColumnX - 20)

	// LEFT COLUMN: Event history
	// Display condensed event history
	pdf.SetFont("AcuminMedium", "", 12) // Increased font size by 4pt (was 8pt)
	eventY := startY

	lineHeight := float64(6) // Increased line height for larger font
	maxEvents := 20 // Maximum events to display

	for i, event := range artist.EventHistory {
		if i >= maxEvents {
			pdf.SetXY(leftColumnX, eventY + float64(i)*lineHeight)
			pdf.Cell(columnWidth, lineHeight, fmt.Sprintf("... and %d more events", len(artist.EventHistory)-maxEvents))
			break
		}

		pdf.SetXY(leftColumnX, eventY + float64(i)*lineHeight)
		// Display event details in a condensed format with winner status
		winnerText := ""
		if event.IsWinner {
			winnerText = " W" // Add W for winners
		}
		pdf.Cell(columnWidth, lineHeight, fmt.Sprintf("%s R%d-E%d%s", event.EventEID, event.Round, event.EaselNumber, winnerText))
	}

	// RIGHT COLUMN: Artist Bio
	pdf.SetFont("AcuminMedium", "", 14) // Increased font size by 4pt (was 10pt)
	pdf.SetXY(rightColumnX, startY)

	if artist.Bio != "" {
		// Use MultiCell for automatic word wrapping
		pdf.MultiCell(rightColumnMaxWidth, 6, cleanString(artist.Bio), "", "L", false)
	} else {
		pdf.Cell(rightColumnMaxWidth, 6, "No bio available")
	}
}

// cleanString removes problematic characters that can cause PDF issues
func cleanString(s string) string {
	// Replace problematic Unicode characters
	s = strings.ReplaceAll(s, "\u201c", "\"")  // left double quotation mark
	s = strings.ReplaceAll(s, "\u201d", "\"")  // right double quotation mark
	s = strings.ReplaceAll(s, "\u2018", "'")   // left single quotation mark
	s = strings.ReplaceAll(s, "\u2019", "'")   // right single quotation mark
	s = strings.ReplaceAll(s, "\u2013", "-")   // en dash
	s = strings.ReplaceAll(s, "\u2014", "-")   // em dash
	s = strings.ReplaceAll(s, "\u2026", "...") // horizontal ellipsis

	// Remove null bytes and other control characters
	s = strings.ReplaceAll(s, "\x00", "")

	return s
}