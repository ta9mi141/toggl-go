package reports

import (
	"context"
	"fmt"
)

const (
	detailedEndpoint string = "/reports/api/v2/details"
)

// DetailedRequestParameters represents request parameters used in the detailed report.
type DetailedRequestParameters struct {
	*StandardRequestParameters
	Page int
}

func (params *DetailedRequestParameters) urlEncode() string {
	values := params.StandardRequestParameters.values()

	if params.Page != 0 {
		values.Add("page", fmt.Sprint(params.Page))
	}

	return values.Encode()
}

func (c *Client) GetDetailed(ctx context.Context, params *DetailedRequestParameters, detailedReport interface{}) error {
	err := c.get(ctx, c.buildURL(detailedEndpoint, params), detailedReport)
	if err != nil {
		return err
	}
	return nil
}
