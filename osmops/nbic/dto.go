package nbic

import (
	"encoding/json"

	"github.com/fluxcd/source-watcher/osmops/nbic/sec"
)

type nbiTokenPayloadView struct {
	Id      string  `json:"id"`
	Expires float64 `json:"expires"`
}

// FromNbiTokenPayload builds a Token out of the JSON payload returned by a
// call to the NBI token endpoint.
func FromNbiTokenPayload(data []byte) (*sec.Token, error) {
	payload := nbiTokenPayloadView{}
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, err
	}
	return sec.NewToken(payload.Id, payload.Expires), nil
}
