package ksyun

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/hashcode"
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
}

func SchemaSetToFilterMap(s interface{}, prefix string, index int, input *map[string]interface{}) {
	(*input)["Filter."+strconv.Itoa(index)+".Name"] = prefix
	count := int(0)
	for _, v := range s.(*schema.Set).List() {
		count = count + 1
		(*input)["Filter."+strconv.Itoa(index)+".Value."+strconv.Itoa(count)] = v
	}
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
	switch data:=data.(type) {
	case string:
		bs = []byte(data)
	default:
		bs, err = json.MarshalIndent(data, "", "\t")
		if err != nil {
			return fmt.Errorf("MarshalIndent data %#v and got an error: %#v", data, err)
		}
	}

	return ioutil.WriteFile(absPath, bs, 0422)
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

// schemaSetToStringSlice used for converting terraform schema set to a string slice
func SchemaSetToStringSlice(s interface{}) []string {
	vL := []string{}

	for _, v := range s.(*schema.Set).List() {
		vL = append(vL, v.(string))
	}

	return vL
}
