package model

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
)

// from https://github.com/99designs/gqlgen/issues/924#issuecomment-558690205

//MarshalInt64 marshall int64
func MarshalInt64(t int64) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = io.WriteString(w, strconv.FormatInt(t, 10))
	})
}

//UnmarshalInt64 unmarshall int64
func UnmarshalInt64(v interface{}) (int64, error) {
	if res, ok := v.(json.Number); ok {
		return res.Int64()
	}
	if res, ok := v.(string); ok {
		return json.Number(res).Int64()
	}
	if res, ok := v.(int64); ok {
		return res, nil
	}
	if res, ok := v.(*int64); ok {
		return *res, nil
	}
	return 0, fmt.Errorf("could not convert %v of type %T to Int64", v, v)
}
