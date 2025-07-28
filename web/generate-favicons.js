import sharp from 'sharp';
import fs from 'fs';
import path from 'path';

// Create public directory if it doesn't exist
const publicDir = './public';
if (!fs.existsSync(publicDir)) {
  fs.mkdirSync(publicDir, { recursive: true });
}

// Create a supermarket-themed favicon SVG with orange colors
const faviconSvg = `
<svg width="512" height="512" viewBox="0 0 512 512" xmlns="http://www.w3.org/2000/svg">
  <!-- Background gradient from orange to darker orange -->
  <defs>
    <linearGradient id="bgGradient" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" style="stop-color:#FF8C00;stop-opacity:1" />
      <stop offset="100%" style="stop-color:#FF4500;stop-opacity:1" />
    </linearGradient>
    <linearGradient id="cartGradient" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" style="stop-color:#FFFFFF;stop-opacity:1" />
      <stop offset="100%" style="stop-color:#F0F0F0;stop-opacity:1" />
    </linearGradient>
  </defs>
  
  <!-- Background -->
  <rect width="512" height="512" fill="url(#bgGradient)" rx="80"/>
  
  <!-- Shopping cart body -->
  <rect x="120" y="200" width="200" height="140" rx="20" fill="url(#cartGradient)" stroke="#333" stroke-width="8"/>
  
  <!-- Shopping cart handle -->
  <rect x="80" y="220" width="60" height="12" rx="6" fill="url(#cartGradient)" stroke="#333" stroke-width="4"/>
  
  <!-- Shopping cart wheels -->
  <circle cx="160" cy="380" r="25" fill="#333"/>
  <circle cx="280" cy="380" r="25" fill="#333"/>
  <circle cx="160" cy="380" r="15" fill="#666"/>
  <circle cx="280" cy="380" r="15" fill="#666"/>
  
  <!-- Products in cart (representing supermarket items) -->
  <rect x="140" y="220" width="30" height="40" rx="5" fill="#FF6B35"/>
  <rect x="180" y="210" width="25" height="50" rx="5" fill="#4CAF50"/>
  <rect x="215" y="225" width="35" height="35" rx="5" fill="#2196F3"/>
  <rect x="260" y="215" width="30" height="45" rx="5" fill="#9C27B0"/>
  
  <!-- Discount percentage symbol -->
  <circle cx="380" cy="150" r="60" fill="#FFD700" stroke="#FF4500" stroke-width="6"/>
  <text x="380" y="140" text-anchor="middle" font-family="Arial Black, sans-serif" font-size="32" font-weight="bold" fill="#FF4500">%</text>
  <text x="380" y="175" text-anchor="middle" font-family="Arial, sans-serif" font-size="24" font-weight="bold" fill="#FF4500">OFF</text>
  
  <!-- Money symbol -->
  <circle cx="420" cy="320" r="40" fill="#32CD32" stroke="#228B22" stroke-width="4"/>
  <text x="420" y="335" text-anchor="middle" font-family="Arial Black, sans-serif" font-size="36" font-weight="bold" fill="white">$</text>
</svg>
`.trim();

// Write the SVG favicon
fs.writeFileSync(path.join(publicDir, 'vite.svg'), faviconSvg);

// Sizes for different favicon formats
const faviconSizes = [
  { size: 16, name: 'favicon-16x16.png' },
  { size: 32, name: 'favicon-32x32.png' },
];

const appleTouchIconSizes = [
  { size: 57, name: 'apple-touch-icon-57x57.png' },
  { size: 60, name: 'apple-touch-icon-60x60.png' },
  { size: 72, name: 'apple-touch-icon-72x72.png' },
  { size: 76, name: 'apple-touch-icon-76x76.png' },
  { size: 114, name: 'apple-touch-icon-114x114.png' },
  { size: 120, name: 'apple-touch-icon-120x120.png' },
  { size: 144, name: 'apple-touch-icon-144x144.png' },
  { size: 152, name: 'apple-touch-icon-152x152.png' },
  { size: 180, name: 'apple-touch-icon.png' }, // Default apple-touch-icon
];

// Generate all favicon sizes
async function generateFavicons() {
  console.log('Generating supermarket-themed favicons...');

  // Create base PNG from SVG
  const basePng = await sharp(Buffer.from(faviconSvg))
    .resize(512, 512)
    .png()
    .toBuffer();

  // Generate standard favicons
  for (const { size, name } of faviconSizes) {
    await sharp(basePng)
      .resize(size, size)
      .png()
      .toFile(path.join(publicDir, name));
    console.log(`Generated ${name}`);
  }

  // Generate Apple Touch Icons
  for (const { size, name } of appleTouchIconSizes) {
    await sharp(basePng)
      .resize(size, size)
      .png()
      .toFile(path.join(publicDir, name));
    console.log(`Generated ${name}`);
  }

  console.log('All supermarket-themed favicons generated successfully!');
}

generateFavicons().catch(console.error);
