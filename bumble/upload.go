package bumble

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"maunium.net/go/mautrix/id"
)

func (bum *BumbleConnector) uploadImageFromURL(ctx context.Context, url string) (id.ContentURI, error) {
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		return id.ContentURI{}, fmt.Errorf("download image: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return id.ContentURI{}, fmt.Errorf("read image: %w", err)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = http.DetectContentType(data)
	}

	uploadResp, err := bum.kit.Bot.UploadBytes(ctx, data, contentType)
	if err != nil {
		return id.ContentURI{}, fmt.Errorf("upload image to Matrix: %w", err)
	}

	return uploadResp.ContentURI, nil
}

