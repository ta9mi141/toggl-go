package reports

import (
	"context"
)

const (
	summaryEndpoint string = "/reports/api/v2/summary"
)

type SummaryRequestParameters struct {
	*StandardRequestParameters
	Grouping            string
	Subgrouping         string
	SubgroupingIds      bool
	GroupedTimeEntryIds bool
}

func (params *SummaryRequestParameters) urlEncode() string {
	values := params.StandardRequestParameters.values()

	if params.Grouping != "" {
		values.Add("grouping", params.Grouping)
	}
	if params.Subgrouping != "" {
		values.Add("subgrouping", params.Subgrouping)
	}
	if params.GroupedTimeEntryIds == true {
		values.Add("grouped_time_entry_ids", "true")
	}
	if params.SubgroupingIds == true {
		values.Add("subgrouping_ids", "true")
	}

	return values.Encode()
}

func (c *client) GetSummary(ctx context.Context, params *SummaryRequestParameters, summaryReport interface{}) error {
	err := c.get(ctx, c.buildURL(summaryEndpoint, params), summaryReport)
	if err != nil {
		return err
	}
	return nil
}
