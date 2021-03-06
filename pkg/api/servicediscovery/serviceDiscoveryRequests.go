package servicediscovery

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/erply/api-go-wrapper/internal/common"
	sharedCommon "github.com/erply/api-go-wrapper/pkg/api/common"
	"github.com/pkg/errors"
)

func (cli *Client) GetServiceEndpoints(ctx context.Context) (*ServiceEndpoints, error) {
	const method = "getServiceEndpoints"
	resp, err := cli.SendRequest(ctx, method, map[string]string{})
	if err != nil {
		return nil, err
	}

	res := &getServiceEndpointsResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to decode %s response", method))
	}
	if !common.IsJSONResponseOK(&res.Status) {
		return nil, sharedCommon.NewFromResponseStatus(&res.Status)
	}
	if len(res.Records) < 1 {
		return nil, errors.New("no records in response")
	}
	return &res.Records[0], nil
}
