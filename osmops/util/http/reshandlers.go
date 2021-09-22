package http

import (
	"fmt"
	"net/http"

	jsoniter "github.com/json-iterator/go"

	"github.com/fluxcd/source-watcher/osmops/util"
)

type jsonResReader struct {
	deserialized        interface{}
	expectedStatusCodes util.IntSet
}

func (r *jsonResReader) Handle(res *http.Response) error {
	if res == nil {
		return fmt.Errorf("nil response")
	}
	if r.deserialized == nil {
		return fmt.Errorf("nil deserialization target")
	}
	if !r.expectedStatusCodes.Contains(res.StatusCode) {
		return fmt.Errorf("unexpected response status: %s", res.Status)
	}

	var json = jsoniter.ConfigCompatibleWithStandardLibrary // (*)
	decoder := json.NewDecoder(res.Body)
	return decoder.Decode(r.deserialized)

	// (*) json-iterator lib.
	// We use it in the JsonBody builder to work around encoding/json's
	// inability to serialise map[interface {}]interface{} types. Here
	// we're parsing JSON into a data structure and AFAICT the built-in
	// json lib can parse pretty much any valid JSON you throw at it.
	// So the only reason to use json-iterator in place of encoding/json
	// is performance: json-iterator is way faster than encoding/json.
}

// ReadJsonResponse builds a ResHandler to deserialise a JSON response body.
// The response is expected to be a 200 with a JSON body and the body gets
// read into the provided data structure target pointer. Optionally, you can
// specify expected response codes if you're expecting something different
// than a 200. ReadJsonResponse will return an error if the response code
// isn't among the ones you specified. Also, ReadJsonResponse returns any
// other error that stopped it from deserializing the response body.
//
// Example.
//
//     client := &http.Client{Timeout: time.Second * 10}
//     target := &MyData{}
//     Request(
//         GET, At(url),
//         Accept(MediaType.JSON),
//     ).
//     SetHandler(ReadJsonResponse(target)).
//     RunWith(client.Do)
//
func ReadJsonResponse(target interface{}, expectedStatusCode ...int) ResHandler {
	if len(expectedStatusCode) == 0 {
		expectedStatusCode = append(expectedStatusCode, http.StatusOK)
	}
	return &jsonResReader{
		deserialized:        target,
		expectedStatusCodes: util.ToIntSet(expectedStatusCode...),
	}
}
