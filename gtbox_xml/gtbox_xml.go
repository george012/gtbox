package gtbox_xml

import (
	"bytes"
	"encoding/xml"
	"io"
)

func GetValueWithKey(xmlData []byte, key string) string {
	decoder := xml.NewDecoder(bytes.NewReader(xmlData))
	var keyValue string
	for {
		t, err := decoder.Token()
		if err == io.EOF {
			break
		}

		if err != nil {
			return ""
		}
		switch token := t.(type) {
		case xml.StartElement:
			if token.Name.Local == key {
				if f_err := decoder.DecodeElement(&keyValue, &token); f_err != nil {
					return ""
				}
				return keyValue
			}
		}
	}
	return ""
}
