package strtool

import "bytes"

func UnderlineString(s string) (res string) {
	A := 'A'
	Z := 'Z'
	a := 'a'
	delta := a - A
	buff := bytes.Buffer{}
	for i, b := range s {
		if b >= A && b <= Z {
			if i == 0 {
				buff.WriteByte(byte(b + delta))
			} else {
				buff.WriteByte('_')
				buff.WriteByte(byte(b + delta))
			}
			continue
		}

		buff.WriteByte(byte(b))
	}

	return buff.String()
}
