package craiyon

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"net/http"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/storage"
)

const API_URL = "https://backend.craiyon.com/generate"

type Client struct {
	client     *http.Client
	gcs        *storage.Client
	bucketName string
}

type Request struct {
	Prompt string `json:"prompt"`
}

type Response struct {
	Images []string `json:"images,omitempty"`
	// TODO figure out error response
}

func NewClient() *Client {
	bn := os.Getenv("BUCKET_NAME")
	if bn == "" {
		panic(fmt.Errorf("BUCKET_NAME cannot be empty"))
	}

	gcs, err := storage.NewClient(context.Background())
	if err != nil {
		panic(err)
	}
	c := &http.Client{
		Timeout: 5 * time.Minute,
	}

	return &Client{client: c, gcs: gcs, bucketName: bn}
}

func (c *Client) GetImage(text string) (string, error) {
	ctx := context.Background()

	body := &Request{Prompt: text}
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(body)
	if err != nil {
		return "", fmt.Errorf("error marshalling body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, API_URL, b)
	if err != nil {
		return "", fmt.Errorf("error creating http request: %w", err)
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("user-agent", "please don't block me i'm using this for a conference talk - @eddiezane")

	res, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error posting: %w", err)
	}
	defer res.Body.Close()

	rj := new(Response)
	err = json.NewDecoder(res.Body).Decode(rj)
	if err != nil {
		return "", fmt.Errorf("error decoding response: %w", err)
	}

	metaImage := image.NewRGBA(image.Rect(0, 0, 256*3, 256*3))

	for i, raw := range rj.Images {
		s := strings.ReplaceAll(raw, `\n`, ``)
		jpg, err := base64.StdEncoding.DecodeString(s)
		if err != nil {
			return "", fmt.Errorf("error decoding base64: %w", err)
		}

		img, err := jpeg.Decode(bytes.NewReader(jpg))
		if err != nil {
			return "", fmt.Errorf("error decoding jpeg: %w", err)
		}

		row := i / 3
		col := i % 3

		startx := col * 256
		starty := row * 256

		draw.Draw(metaImage, image.Rect(startx, starty, startx+256, starty+256), img, image.Point{}, draw.Src)
	}

	sum := md5.Sum([]byte(text))

	fileName := fmt.Sprintf("%s-%d.jpeg", hex.EncodeToString(sum[:]), time.Now().Unix())

	oh := c.gcs.Bucket(c.bucketName).Object(fileName)

	writer := oh.NewWriter(ctx)

	err = jpeg.Encode(writer, metaImage, &jpeg.Options{Quality: 100})
	if err != nil {
		return "", fmt.Errorf("error writing file: %w", err)
	}
	if err = writer.Close(); err != nil {
    return "", fmt.Errorf("error writing file to cloud: %w", err)
  }

	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", c.bucketName, fileName), nil
}
