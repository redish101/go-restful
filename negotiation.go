package restful

import (
	"sort"
	"strconv"
	"strings"
)

func negotiate(accept string, candidates []string) string {
	if len(candidates) == 0 {
		return ""
	}
	if accept == "" || accept == "*/*" {
		return candidates[0]
	}

	type pref struct {
		mime string
		q    float64
		idx  int
	}

	var prefs []pref
	for i, part := range strings.Split(accept, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		seg := strings.Split(part, ";")
		mime := strings.TrimSpace(seg[0])
		q := 1.0
		for _, p := range seg[1:] {
			p = strings.TrimSpace(p)
			if strings.HasPrefix(p, "q=") {
				if v, err := strconv.ParseFloat(p[2:], 64); err == nil {
					q = v
				}
			}
		}
		prefs = append(prefs, pref{mime: mime, q: q, idx: i})
	}

	sort.Slice(prefs, func(i, j int) bool {
		if prefs[i].q != prefs[j].q {
			return prefs[i].q > prefs[j].q
		}
		return prefs[i].idx < prefs[j].idx
	})

	for _, p := range prefs {
		for _, cand := range candidates {
			if matchMIME(p.mime, cand) {
				return cand
			}
		}
	}
	return ""
}

func matchMIME(pattern, mime string) bool {
	if pattern == "*/*" || pattern == mime {
		return true
	}
	if strings.HasSuffix(pattern, "/*") {
		prefix := strings.TrimSuffix(pattern, "/*")
		return strings.HasPrefix(mime, prefix+"/")
	}
	return false
}
