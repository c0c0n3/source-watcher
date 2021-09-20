package nbic

import (
	"fmt"
)

type nsDescView struct { // only the response fields we care about.
	Id   string `json:"_id"`
	Name string `json:"id"`
}

type nsDescMap map[string]string

func buildNsDescMap(ds []nsDescView) nsDescMap {
	descMap := map[string]string{}
	for _, d := range ds {
		descMap[d.Name] = d.Id
	}
	return descMap
}

func (c *Session) getNsDescriptors() ([]nsDescView, error) {
	data := []nsDescView{}
	if _, err := c.getJson(c.conn.NsDescriptors(), &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Session) lookupNsDescriptorId(name string) (string, error) {
	if c.nsdMap == nil {
		if ds, err := c.getNsDescriptors(); err != nil {
			return "", err
		} else {
			c.nsdMap = buildNsDescMap(ds)
		}
	}
	if id, ok := c.nsdMap[name]; !ok {
		return "", fmt.Errorf("no nsd found for name ID: %s", name)
	} else {
		return id, nil
	}
}
