package data

import "strings"

type BlackListRule struct {
	// blackListType: endsWidth or absolute. the reason it is not a enum is so that
	// sdk does not break in case there is a version mismatch between sdk and
	// commons
	BlackListType string            `json:"blackListType"`
	MatchValues   []string          `json:"matchValues"`
	Method        HTTPRequestMethod `json:"method"`
}

func (b *BlackListRule) IsValid() bool {
	if b.BlackListType == "" {
		return false
	}
	if b.MatchValues == nil {
		return false
	}
	return true
}

func (b *BlackListRule) MatchUri(uri URI, method HTTPRequestMethod) bool {
	if !b.IsValid() {
		return false
	}
	if len(b.Method) > 0 && b.Method != method {
		return false
	}
	if b.BlackListType == "endsWith" {
		for _, matchValue := range b.MatchValues {
			if strings.HasSuffix(strings.ToLower(uri.GetURIPath()), strings.ToLower(matchValue)) {
				return true
			}
		}
	} else if b.BlackListType == "absolute" {
		for _, matchValue := range b.MatchValues {
			if uri.Equals(GetURI(matchValue)) {
				return true
			}
		}
	}
	if strings.EqualFold(b.BlackListType, "endsWith") {
		for _, matchValue := range b.MatchValues {
			if strings.HasSuffix(strings.ToLower(uri.UriPath), strings.ToLower(matchValue)) {
				return true
			}
		}
	} else if strings.EqualFold(b.BlackListType, "absolute") {
		for _, matchValue := range b.MatchValues {
			if uri.String() == matchValue {
				return true
			}
		}
	}
	return false
}
