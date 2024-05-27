<p align="center" width="100%">
    <img width="40%" style="image-rendering: pixelated;" src="summer_Illumine.gif"> 
</p>

<p align="center" style="font-size: 2rem;">GoTerPix</p>
<p align="center" style="font-size: 1.5rem;">Pixel art images right into your terminal!</p>

# Build
    git clone https://github.com/FinecoFinit/goterpix.git
    cd goterpix
    go build goterpix.go

# Usage
    ./goterpix summer_Illumine.gif 60

goterpix | path to gif image | frame delay in milliseconds

## Limitations
- Only gif images are usable at the moment
- Image should be perfect 1 to 1 pixels 