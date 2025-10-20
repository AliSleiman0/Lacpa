# Adding Firm Logos - Quick Guide

This guide explains how to add or update firm logos in the LACPA system.

## Directory Structure

```
LACPA_Web/
  assets/
    images/
      firms/
        deloitte.svg
        pwc.svg
        ey.svg
        kpmg.svg
        bdo.svg
        grant-thornton.svg
        semaan-gholam.png
        fiduciaire.png
        khalil.png
        najjar.png
```

## Logo Specifications

### Recommended Format
- **SVG** (vector) - Best for scalability and file size
- **PNG** - Good for complex logos with transparency
- **JPG** - Only for photos or complex images without transparency

### Recommended Dimensions
- **Width:** 200px
- **Height:** 80px
- **Aspect Ratio:** 2.5:1 (landscape orientation)
- **Background:** Transparent (PNG/SVG) or solid color

### File Size
- Try to keep under 50KB per logo
- SVG files are typically 2-10KB
- Optimize PNGs using tools like TinyPNG

## How to Add a New Logo

### Step 1: Prepare the Logo File

1. Obtain the firm's logo (from their website or brand guidelines)
2. Resize to recommended dimensions (200x80px)
3. Save in appropriate format (SVG preferred)
4. Name the file clearly (e.g., `company-name.svg`)

### Step 2: Add to Assets Folder

```bash
# Copy logo to the firms folder
cp your-logo.svg LACPA_Web/assets/images/firms/
```

### Step 3: Update Seed Script

Edit `Backend/scripts/seed_firms.js`:

```javascript
{
    _id: new ObjectId(),
    lacpa_id: "FIRM-2024-XXX",
    firm_name: "Your Firm Name",
    logo_url: "/assets/images/firms/your-logo.svg",  // ← Add this line
    firm_type: "Audit Firm",
    // ... rest of the fields
}
```

### Step 4: Re-seed Database

```bash
cd Backend
node scripts/seed_firms.js
```

### Step 5: Verify

Visit http://localhost:3000/membership/firms and check that the logo displays correctly.

## Logo Display Behavior

### On the Frontend
- Logos are displayed in a **white container** (128px × 96px)
- Images are **scaled proportionally** to fit
- **White background** helps logos with dark backgrounds stand out
- **Rounded corners** (8px border-radius) for modern look

### Fallback
If no logo is provided (`logo_url` is empty):
- A **building emoji** (🏢) is shown
- Gray background for consistency

## Creating Simple SVG Logos

If you need to create a basic text logo, here's a template:

```svg
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 200 80">
  <!-- Background -->
  <rect width="200" height="80" fill="#0076A8"/>
  
  <!-- Company Name -->
  <text 
    x="100" 
    y="50" 
    font-family="Arial, sans-serif" 
    font-size="28" 
    font-weight="bold" 
    fill="white" 
    text-anchor="middle">
    Your Company
  </text>
</svg>
```

### Customization
- Change `fill="#0076A8"` to your brand color
- Adjust `font-size="28"` for text size
- Modify `y="50"` to vertically center text
- For two-line logos, use multiple `<text>` elements

## Troubleshooting

### Logo Not Showing
1. Check file path is correct: `/assets/images/firms/filename.ext`
2. Verify file exists in `LACPA_Web/assets/images/firms/`
3. Check browser console for 404 errors
4. Clear browser cache (Ctrl+Shift+R)

### Logo Too Small/Large
1. Open the SVG in a text editor
2. Adjust `viewBox="0 0 200 80"` dimensions
3. Or use `width` and `height` attributes

### Logo Colors Wrong
- Ensure logo has transparent background (PNG/SVG)
- Or ensure background color complements the white container
- For dark logos, consider adding a light version

### Logo Blurry
- Use SVG instead of PNG/JPG for sharp rendering
- If using PNG, ensure it's at least 200x80px
- Use 2x resolution (400x160px) for retina displays

## Best Practices

1. **Use SVG when possible** - Scales perfectly, small file size
2. **Transparent backgrounds** - Let the container color show through
3. **Simple is better** - Logos should be recognizable at small sizes
4. **Consistent style** - Keep similar padding/margins across all logos
5. **Test on different screens** - Check mobile and desktop views
6. **Brand guidelines** - Follow the firm's official brand guidelines

## Batch Processing

To add multiple logos at once:

```bash
# Copy all logos
cp logos/*.svg LACPA_Web/assets/images/firms/

# Update seed_firms.js with all new entries
# Re-seed once
node scripts/seed_firms.js
```

## External Resources

- **Logo Download:** Company websites, Wikipedia, Wikimedia Commons
- **SVG Editing:** [Figma](https://figma.com), [Inkscape](https://inkscape.org)
- **PNG Optimization:** [TinyPNG](https://tinypng.com)
- **SVG Optimization:** [SVGOMG](https://jakearchibald.github.io/svgomg/)
