package pdf_export

// Response — the data sent to the frontend after a successful export. The
// frontend can use Path to show/open the file (e.g. via Wails runtime's
// file-explorer helpers) without needing to know how the export directory
// was resolved.
type Response struct {
	FileName string `json:"fileName"`
	Path     string `json:"path"`
}
