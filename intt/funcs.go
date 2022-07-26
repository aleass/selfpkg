package intt

import "fmt"

func NumChanges(v float64) string {
	if v < 1e8 {
		v = v / 1e4
		return fmt.Sprintf("%0.2f%s", v, "万")
	} else {
		v = v / 1e8
		return fmt.Sprintf("%0.2f%s", v, "亿")
	}

}
