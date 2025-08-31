# ChatGPT Icon Generation Prompt for ImageKit

## Main Prompt for GPT-5 (DALL-E 3)

```
Create a modern, minimalist app icon for "ImageKit" - a privacy-focused web-based image conversion tool.

Design Requirements:
- Square format, 512x512 pixels
- **TRANSPARENT BACKGROUND** (PNG format with alpha channel)
- Clean, flat design with subtle depth
- Must work on both light and dark backgrounds
- Professional app store quality

Visual Concept:
Create an abstract representation combining:
1. A stylized image/photo frame or layers (representing image processing)
2. Transformation arrows or geometric shapes showing conversion/change
3. Use a purple gradient (#667eea to #764ba2) as the primary color
4. White or light accents for contrast

Style Guidelines:
- Modern, minimalist aesthetic
- No text or letters
- Geometric and clean shapes
- Subtle shadows or gradients for depth
- Should be recognizable at small sizes (16x16)

The icon should convey: transformation, privacy, speed, and professionalism.
Make it look like a premium app icon you'd see in Apple App Store or Google Play Store.
```

## Alternative Prompt (More Specific)

```
Design a 512x512 pixel app icon with these exact specifications:

Background: TRANSPARENT (PNG with alpha channel)
Icon Shape: Rounded square with purple gradient (#667eea top-left to #764ba2 bottom-right)

Main Element: Create a white, minimalist symbol in the center that combines:
- Two overlapping rectangles (representing before/after images)
- A subtle arrow or transformation indicator between them
- Keep it abstract and geometric

Style: Similar to modern app icons like Notion, Figma, or Slack - clean, professional, with subtle depth through gradients.

The icon must be visually balanced, work at all sizes from 16x16 to 512x512, and convey image transformation without being literal.
```

## After Generation

Once you receive the icon from ChatGPT:

1. Save the main 512x512 icon
2. Use the provided script to generate all sizes:
   ```bash
   ./scripts/generate-all-icon-sizes.sh path/to/icon-512x512.png
   ```

3. Or manually create sizes using an online tool like:
   - https://www.favicon-generator.org/
   - https://realfavicongenerator.net/
   - https://www.pwabuilder.com/imageGenerator

## Required Icon Sizes

- 16x16 (favicon)
- 32x32 (favicon)
- 72x72 (Android)
- 96x96 (Android)
- 128x128 (Chrome Web Store)
- 144x144 (Android)
- 152x152 (iOS)
- 180x180 (iOS)
- 192x192 (Android/PWA)
- 384x384 (Android/PWA)
- 512x512 (Android/PWA)

## Color Reference

- Primary: #667eea
- Secondary: #764ba2
- Accent: #f093fb
- Background (optional): #ffffff or transparent