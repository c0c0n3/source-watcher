package nbic

import "fmt"

type nsInstanceView struct { // only the response fields we care about.
	Id   string `json:"_id"`
	Name string `json:"name"`
}

type nsInstanceMap map[string][]string

func (m nsInstanceMap) addMapping(name string, id string) {
	if entry, ok := m[name]; ok {
		m[name] = append(entry, id)
	} else {
		m[name] = []string{id}
	}
}

func buildNsInstanceMap(vs []nsInstanceView) nsInstanceMap {
	nsMap := nsInstanceMap{}
	for _, v := range vs {
		nsMap.addMapping(v.Name, v.Id)
	}
	return nsMap
}

// NOTE. NS instance name to ID lookup.
// OSM NBI doesn't enforce uniqueness of NS names. In fact, it lets you happily
// create a new instance even if an existing one has the same name, e.g.
//
// $ curl localhost/osm/nslcm/v1/ns_instances_content \
// -v -X POST \
// -H 'Authorization: Bearer 0WhgBufy1Wt82NbF9OsmftwpRfcsV4sU' \
// -H 'Content-Type: application/yaml' \
// -d'{"nsdId": "aba58e40-d65f-4f4e-be0a-e248c14d3e03", "nsName": "ldap", "nsDescription": "default description", "vimAccountId": "4a4425f7-3e72-4d45-a4ec-4241186f3547"}'
// ...
// HTTP/1.1 201 Created
// ...
// ---
// id: 794ef9a2-8bbb-42c1-869a-bab6422982ec
// nslcmop_id: 0fdfaa6a-b742-480c-9701-122b3f732e4
//
// This is why we map an NS instance name to a list of IDs.

func (c *Session) getNsInstancesContent() ([]nsInstanceView, error) {
	data := []nsInstanceView{}
	if _, err := c.getJson(c.conn.NsInstancesContent(), &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Session) lookupNsInstanceId(name string) (string, error) {
	if c.nsInstMap == nil {
		if vs, err := c.getNsInstancesContent(); err != nil {
			return "", err
		} else {
			c.nsInstMap = buildNsInstanceMap(vs)
		}
	}
	if ids, ok := c.nsInstMap[name]; !ok {
		return "", fmt.Errorf("no NS instance found for name: %s", name)
	} else {
		if len(ids) != 1 {
			return "",
				fmt.Errorf("NS instance name not bound to a single ID: %v", ids)
		}
		return ids[0], nil
	}
}
