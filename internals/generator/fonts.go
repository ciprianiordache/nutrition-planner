package generator

import _ "embed"

// Default fonts embedded directly into the binary, so pdfgen works out of
// the box with zero external files. A GeneratorConfig can still override
// FontRegularTTF/FontBoldTTF with a path to use a different typeface -
// embedding here is just the zero-config fallback.

//go:embed assets/fonts/NotoSans-Regular.ttf
var defaultFontRegular []byte

//go:embed assets/fonts/NotoSans-Bold.ttf
var defaultFontBold []byte
