package views

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strings"
)

// https://blog.questionable.services/article/approximating-html-template-inheritance/
// https://github.com/asit-dhal/golang-template-layout

var FuncMap = template.FuncMap{
	"bytes":       Bytes,
	"icon":        Icon,
	"marshalJSON": marshalJSONScript}

func Bytes(size int64) string {
	if size < 1000 {
		return fmt.Sprintf("%d Bytes", size)
	}
	size2 := float64(size)
	//ext := []string{"B", "KiB", "MiB", "GiB"}
	//i := 0
	//for ; size > 1024; i++ {
	//	size = size / 1024
	//}
	ext2 := []string{"B", "kB", "MB", "GB"}
	j := 0
	for ; size2 > 1000; j++ {
		size2 = size2 / 1000.0
	}
	//fmt.Printf("%d %s (%.2f %s)", size, ext[i], size2, ext2[j])
	return fmt.Sprintf("%.1f %s", size2, ext2[j])
}
func Icon(typ string) string {
	ext := map[string]string{"F": "file", "FI": "image", "FT": "file-text", "D": "folder", "DA": "album"}
	if icon, ok := ext[typ]; ok {
		return icon
	}
	return ""
}

func join(strs ...string) string {
	var sb strings.Builder
	for _, str := range strs {
		sb.WriteString(str)
	}
	return sb.String()
}

func marshalJSONScript(name string, v interface{}) template.HTML {

	startTag := "<script id='" + name + "' type='application/json'>"
	endTag := "</script>"
	json, err := json.Marshal(v)
	//fmt.Println(string(json))
	if err != nil {
		return template.HTML("Error marschalling JSON: " + err.Error())
	}
	script := startTag + string(json) + endTag
	return template.HTML(script)
}
