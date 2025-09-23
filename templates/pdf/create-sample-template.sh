#!/bin/bash

# This script creates a basic PDF template for testing
# Requires: imagemagick and ghostscript

echo "Creating sample PDF templates..."

# Create artist page template
convert -size 2790x2160 xc:white \
  -fill '#2980B9' -draw 'rectangle 0,0 2790,300' \
  -fill '#E74C3C' -draw 'rectangle 0,300 2790,330' \
  -fill white -pointsize 200 -gravity North -annotate +0+50 'ART BATTLE' \
  -fill white -pointsize 100 -gravity North -annotate +0+180 '[EVENT NAME]' \
  -fill '#CCCCCC' -pointsize 60 -gravity NorthWest -annotate +200+500 '[ARTIST NAME HERE]' \
  -fill '#CCCCCC' -draw 'rectangle 200,800 800,1400' \
  -fill black -pointsize 40 -gravity NorthWest -annotate +220+950 'QR CODE' \
  -fill '#CCCCCC' -pointsize 40 -gravity NorthWest -annotate +1050+900 '[ROUND/EASEL BADGE]' \
  -fill '#CCCCCC' -pointsize 40 -gravity NorthWest -annotate +1050+1000 '[INSTAGRAM]' \
  -fill '#CCCCCC' -pointsize 40 -gravity NorthWest -annotate +1050+1100 '[LOCATION]' \
  -fill '#34495E' -draw 'rectangle 200,1600 2400,1700' \
  -fill white -pointsize 50 -gravity NorthWest -annotate +220+1630 'ARTIST EVENT HISTORY' \
  -fill '#CCCCCC' -draw 'rectangle 200,1700 2400,2000' \
  -density 72 templates/artist-page.pdf

# Create artist list template
convert -size 2790x2160 xc:white \
  -fill '#2980B9' -draw 'rectangle 0,0 2790,400' \
  -fill white -pointsize 200 -gravity North -annotate +0+50 'ART BATTLE' \
  -fill white -pointsize 100 -gravity North -annotate +0+180 '[EVENT NAME]' \
  -fill white -pointsize 60 -gravity North -annotate +0+280 'EVENT CODE: [EID]' \
  -fill '#E74C3C' -draw 'rectangle 0,400 2790,430' \
  -fill '#34495E' -draw 'rectangle 100,500 2690,600' \
  -fill white -pointsize 60 -gravity NorthWest -annotate +120+530 'ARTIST LINEUP' \
  -fill '#CCCCCC' -draw 'rectangle 100,700 2690,1800' \
  -fill black -pointsize 40 -gravity NorthWest -annotate +400+1200 '[ARTIST TABLE WILL BE GENERATED HERE]' \
  -density 72 templates/artist-list.pdf

# Create auction info template
convert -size 2790x2160 xc:white \
  -fill '#2980B9' -draw 'rectangle 0,0 2790,400' \
  -fill white -pointsize 200 -gravity North -annotate +0+50 'ART BATTLE' \
  -fill white -pointsize 100 -gravity North -annotate +0+180 '[EVENT NAME]' \
  -fill white -pointsize 60 -gravity North -annotate +0+280 'EVENT CODE: [EID]' \
  -fill '#E74C3C' -draw 'rectangle 0,400 2790,430' \
  -fill '#34495E' -draw 'rectangle 100,500 2690,600' \
  -fill white -pointsize 60 -gravity NorthWest -annotate +120+530 'AUCTION INFORMATION' \
  -fill '#CCCCCC' -draw 'rectangle 100,700 2690,1800' \
  -fill black -pointsize 40 -gravity NorthWest -annotate +400+1200 '[AUCTION TABLE WILL BE GENERATED HERE]' \
  -density 72 templates/auction-info.pdf

echo "Sample templates created!"
echo "You can now:"
echo "1. Replace these with professional designs from InDesign/Illustrator"
echo "2. Test PDF generation with USE_TEMPLATE_PDF=true"
echo "3. Adjust coordinates in configs/template-config.json as needed"