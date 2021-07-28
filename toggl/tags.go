package toggl

import (
	"context"
	"errors"
	"strconv"
)

const (
	tagsEndpoint string = "/api/v8/tags"
)

var (
	// ErrTagNotFound is returned when the provided tag is nil.
	ErrTagNotFound = errors.New("the provided tag must be non-nil")
)

// Tag represents properties of tag.
type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	WID  int    `json:"wid"`
}

type rawTagData struct {
	Tag Tag `json:"data"`
}

// CreateTag creates a tag.
func (c *Client) CreateTag(ctx context.Context, tag *Tag) (*Tag, error) {
	if tag == nil {
		return nil, ErrTagNotFound
	}
	rawTagData := new(rawTagData)
	if err := c.httpPost(ctx, c.buildURL(tagsEndpoint), tag, rawTagData); err != nil {
		return nil, err
	}
	return &rawTagData.Tag, nil
}

// UpdateTag updates a tag.
func (c *Client) UpdateTag(ctx context.Context, tag *Tag) (*Tag, error) {
	if tag == nil {
		return nil, ErrTagNotFound
	}
	rawTagData := new(rawTagData)
	endpoint := tagsEndpoint + "/" + strconv.Itoa(tag.ID)
	if err := c.httpPut(ctx, c.buildURL(endpoint), tag, rawTagData); err != nil {
		return nil, err
	}
	return &rawTagData.Tag, nil
}

// DeleteTag deletes a tag.
func (c *Client) DeleteTag(ctx context.Context, tag *Tag) error {
	if tag == nil {
		return ErrTagNotFound
	}
	endpoint := tagsEndpoint + "/" + strconv.Itoa(tag.ID)
	if err := c.httpDelete(ctx, c.buildURL(endpoint)); err != nil {
		return err
	}
	return nil
}
