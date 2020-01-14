package ksyun

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

func SchemaSetToInstanceMap(s interface{}, prefix string, input *map[string]interface{}) {

	count := int(0)
	for _, v := range s.(*schema.Set).List() {
		count = count + 1
		(*input)[prefix+"."+strconv.Itoa(count)] = v
	}
	return
}

func SchemaSetToFilterMap(s interface{}, prefix string, index int, input *map[string]interface{}) {

	(*input)["Filter."+strconv.Itoa(index)+".Name"] = prefix
	count := int(0)
	for _, v := range s.(*schema.Set).List() {
		count = count + 1
		(*input)["Filter."+strconv.Itoa(index)+".Value."+strconv.Itoa(count)] = v
	}
	return
}

func SchemaSetsToFilterMap(d *schema.ResourceData, filters []string, req *map[string]interface{}) *map[string]interface{} {
	index := 0
	for _, v := range filters {
		var idsString []string
		if ids, ok := d.GetOk(v); ok {
			idsString = SchemaSetToStringSlice(ids)
		}
		if len(idsString) > 0 {
			index++
			(*req)[fmt.Sprintf("Filter.%v.Name", index)] = strings.Replace(v, "_", "-", -1)
		}
		for k1, v1 := range idsString {
			if v1 == "" {
				continue
			}
			(*req)[fmt.Sprintf("Filter.%v.Value.%d", index, k1+1)] = v1
		}
	}
	return req
}
func hashStringArray(arr []string) string {
	var buf bytes.Buffer

	for _, s := range arr {
		buf.WriteString(fmt.Sprintf("%s-", s))
	}

	return fmt.Sprintf("%d", hashcode.String(buf.String()))
}

func writeToFile(filePath string, data interface{}) error {
	absPath, err := getAbsPath(filePath)
	if err != nil {
		return err
	}

	os.Remove(absPath)

	var bs []byte
	switch data.(type) {
	case string:
		bs = []byte(data.(string))
	default:
		bs, err = json.MarshalIndent(data, "", "\t")
		if err != nil {
			return fmt.Errorf("MarshalIndent data %#v and got an error: %#v", data, err)
		}
	}

	ioutil.WriteFile(absPath, bs, 422)
	return nil
}

func getAbsPath(filePath string) (string, error) {
	if strings.HasPrefix(filePath, "~") {
		usr, err := user.Current()
		if err != nil {
			return "", fmt.Errorf("get current user got an error: %#v", err)
		}

		if usr.HomeDir != "" {
			filePath = strings.Replace(filePath, "~", usr.HomeDir, 1)
		}
	}
	return filepath.Abs(filePath)
}

func merageResultDirect(result *[]map[string]interface{}, source []interface{}) {
	for _, v := range source {
		*result = append(*result, v.(map[string]interface{}))
	}
}

// schemaListToStringSlice used for converting terraform attribute of TypeString embedded in TypeList to a string slice.
// it expected interface{} type as []interface{}, usually get the value from `d.Get` of terraform resource data.
func SchemaListToStringSlice(iface interface{}) []string {
	s := []string{}

	for _, i := range iface.([]interface{}) {
		s = append(s, i.(string))
	}

	return s
}

// schemaSetToStringSlice used for converting terraform schema set to a string slice
func SchemaSetToStringSlice(s interface{}) []string {
	vL := []string{}

	for _, v := range s.(*schema.Set).List() {
		vL = append(vL, v.(string))
	}

	return vL
}

func isResourceTimeoutError(err error) bool {
	timeoutErr, ok := err.(*resource.TimeoutError)
	return ok && timeoutErr.LastError == nil
}
