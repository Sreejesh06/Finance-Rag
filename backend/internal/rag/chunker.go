package rag

import (
	"strings"
	"supplychain-rag/internal/utils"
)

type Chunk struct {
	Text       string
	PageNumber int
	ChunkIndex int
}

func ChunkText(pages []utils.PageContent, chunkSize, chunkOverlap int) []Chunk {
	var chunks []Chunk
	chunkIndex := 0

	for _, page := range pages {
		text := strings.TrimSpace(page.Text)
		if len(text) == 0 {
			continue
		}

		words := strings.Fields(text)

		if len(words) == 0 {
			continue
		}

		// simple word-based chunking
		for i := 0; i < len(words); {
			end := i + chunkSize
			if end > len(words) {
				end = len(words)
			}

			chunkText := strings.Join(words[i:end], " ")
			chunks = append(chunks, Chunk{
				Text:       chunkText,
				PageNumber: page.PageNumber,
				ChunkIndex: chunkIndex,
			})
			chunkIndex++

			if end == len(words) {
				break
			}

			i += chunkSize - chunkOverlap
			if i >= len(words) {
				break
			}
		}
	}

	return chunks
}
