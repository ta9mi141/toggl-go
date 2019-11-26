package toggl

import (
	"context"
	"strconv"
)

const (
	tagsEndpoint string = "/api/v8/tags"
)

// Tag represents properties of tag.
type Tag struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Wid  int    `json:"wid"`
}

// CreateTag creates a tag.
func (c *Client) CreateTag(ctx context.Context, tag *Tag) (*Tag, error) {
	createdTag := new(Tag)
	if err := c.httpPost(ctx, c.buildURL(tagsEndpoint), tag, createdTag); err != nil {
		return nil, err
	}
	return createdTag, nil
}

// UpdateTag updates a tag.
func (c *Client) UpdateTag(ctx context.Context, tag *Tag) (*Tag, error) {
	updatedTag := new(Tag)
	endpoint := tagsEndpoint + "/" + strconv.Itoa(tag.Id)
	if err := c.httpPut(ctx, c.buildURL(endpoint), tag, updatedTag); err != nil {
		return nil, err
	}
	return updatedTag, nil
}

// DeleteTag deletes a tag.
func (c *Client) DeleteTag(ctx context.Context, tag *Tag) error {
	endpoint := tagsEndpoint + "/" + strconv.Itoa(tag.Id)
	if err := c.httpDelete(ctx, c.buildURL(endpoint)); err != nil {
		return err
	}
	return nil
}
