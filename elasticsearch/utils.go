package saelastic

import (
	"encoding/json"
	"log"
	"strings"
	"time"
)

//EventNotification ....
type EventNotification struct {
	Labels      func() interface{}
	Annotations func() interface{}
	StartsAt    string
}

//Sanitize ... Issue with json format in events.
func Sanitize(jsondata string) string {
	r := strings.NewReplacer("\\\"", "\"",
		"\"ves\":\"{", "\"ves\":{",
		"}\"}", "}}",
		"[", "",
		"]", "")
	result := r.Replace(jsondata)
	result1 := strings.Replace(result, "\\\"", "\"", -1)
	return result1

}

//GetIndexNameType ...
func GetIndexNameType(jsondata string) (IndexName, IndexType, error) {
	start := time.Now()

	var f []interface{}
	err := json.Unmarshal([]byte(jsondata), &f)
	if err != nil {
		log.Fatal(err)
		return GENERICINDEX, GENERICINDEXTYPE, err
	}
	elapsed := time.Since(start)
	index, indextype, error := typeSwitch(f[0])
	log.Printf("getIndexNameType took %s", elapsed)
	return index, indextype, error

}

func typeSwitch(tst interface{}) (IndexName, IndexType, error) {
	switch v := tst.(type) {
	case map[string]interface{}:
		if val, ok := v["labels"]; ok {
			switch val.(type) {
			case map[string]interface{}:
				if rec, ok := val.(map[string]interface{}); ok {
					if _, ok := rec["connectivity"]; ok {
						return CONNECTIVITYINDEX, CONNECTIVITYINDEXTYPE, nil
					} else if _, ok := rec["procevent"]; ok {
						return PROCEVENTINDEX, PROCEVENTINDEXTYPE, nil
					} else if _, ok := rec["sysevent"]; ok {
						return SYSEVENTINDEX, SYSEVENTINDEXTYPE, nil
					}
					//else
					return GENERICINDEX, GENERICINDEXTYPE, nil
				} //else
				return GENERICINDEX, GENERICINDEXTYPE, nil
			}
		}
	default:
		return GENERICINDEX, GENERICINDEXTYPE, nil
	}
	return GENERICINDEX, GENERICINDEXTYPE, nil
}