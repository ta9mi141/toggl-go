package reports

import (
	"fmt"
)

const (
	detailedEndpoint string = "/reports/api/v2/details"
)

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

func (c *client) GetDetailed(params *DetailedRequestParameters, detailedReport interface{}) error {
	err := c.get(c.buildURL(detailedEndpoint, params), detailedReport)
	if err != nil {
		return err
	}
	return nil
}
