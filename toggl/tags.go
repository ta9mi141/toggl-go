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
	ErrTagNotFound = errors.New("The provided tag must be non-nil")
)

// Tag represents properties of tag.
type Tag struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Wid  int    `json:"wid"`
}

type rawResponse struct {
	Tag Tag `json:"data"`
}

// CreateTag creates a tag.
func (c *Client) CreateTag(ctx context.Context, tag *Tag) (*Tag, error) {
	if tag == nil {
		return nil, ErrTagNotFound
	}
	rawResponse := new(rawResponse)
	if err := c.httpPost(ctx, c.buildURL(tagsEndpoint), tag, rawResponse); err != nil {
		return nil, err
	}
	return &rawResponse.Tag, nil
}

// UpdateTag updates a tag.
func (c *Client) UpdateTag(ctx context.Context, tag *Tag) (*Tag, error) {
	if tag == nil {
		return nil, ErrTagNotFound
	}
	updatedTag := new(Tag)
	endpoint := tagsEndpoint + "/" + strconv.Itoa(tag.Id)
	if err := c.httpPut(ctx, c.buildURL(endpoint), tag, updatedTag); err != nil {
		return nil, err
	}
	return updatedTag, nil
}

// DeleteTag deletes a tag.
func (c *Client) DeleteTag(ctx context.Context, tag *Tag) error {
	if tag == nil {
		return ErrTagNotFound
	}
	endpoint := tagsEndpoint + "/" + strconv.Itoa(tag.Id)
	if err := c.httpDelete(ctx, c.buildURL(endpoint)); err != nil {
		return err
	}
	return nil
}
