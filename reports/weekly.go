package reports

import (
	"context"
)

const (
	weeklyEndpoint string = "/reports/api/v2/weekly"
)

type WeeklyRequestParameters struct {
	*StandardRequestParameters
	Grouping  string
	Calculate string
}

func (params *WeeklyRequestParameters) urlEncode() string {
	values := params.StandardRequestParameters.values()

	if params.Grouping != "" {
		values.Add("grouping", params.Grouping)
	}
	if params.Calculate != "" {
		values.Add("calculate", params.Calculate)
	}

	return values.Encode()
}

func (c *client) GetWeekly(ctx context.Context, params *WeeklyRequestParameters, weeklyReport interface{}) error {
	err := c.get(ctx, c.buildURL(weeklyEndpoint, params), weeklyReport)
	if err != nil {
		return err
	}
	return nil
}
