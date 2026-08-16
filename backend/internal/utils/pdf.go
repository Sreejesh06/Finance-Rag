package utils

import (
	"bytes"
	"fmt"

	"github.com/ledongthuc/pdf"
)

type PageContent struct {
	PageNumber int
	Text       string
}

func ExtractTextFromPDF(filePath string) ([]PageContent, error) {
	f, r, err := pdf.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open PDF: %w", err)
	}
	defer f.Close()

	var pages []PageContent
	totalPage := r.NumPage()

	for i := 1; i <= totalPage; i++ {
		p := r.Page(i)
		if p.V.IsNull() {
			continue
		}

		var buf bytes.Buffer
		text, err := p.GetPlainText(nil)
		if err != nil {
			// just skip if there is an issue getting text
			continue
		}
		buf.WriteString(text)

		pages = append(pages, PageContent{
			PageNumber: i,
			Text:       buf.String(),
		})
	}

	return pages, nil
}
