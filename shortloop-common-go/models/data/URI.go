package data

import (
	"fmt"
	"strings"
)

type URI struct {
	UriPath         string `json:"uriPath"`
	HasPathVariable bool   `json:"hasPathVariable"`
}

type SegmentTemplateDataType string

const (
	NUMBER  SegmentTemplateDataType = "{number}"
	STRING  SegmentTemplateDataType = "{string}"
	UUID    SegmentTemplateDataType = "{uuid}"
	UNKNOWN SegmentTemplateDataType = "{unknown}"
)

func (s SegmentTemplateDataType) GetPathDisplayName() string {
	switch s {
	case NUMBER:
		return "{number}"
	case STRING:
		return "{string}"
	case UUID:
		return "{uuid}"
	default:
		return "{unknown}"
	}
}

func GetTemplateDataTypeByDisplayName(pathDisplayName string) SegmentTemplateDataType {
	switch pathDisplayName {
	case "{number}":
		return NUMBER
	case "{string}":
		return STRING
	case "{uuid}":
		return UUID
	default:
		return UNKNOWN
	}
}

type URIPathVariable struct {
	variableID   string
	variableName string
}

func (u URIPathVariable) GetVariableID() string {
	return u.variableID
}

func (u URIPathVariable) GetVariableName() string {
	return u.variableName
}

func (u URIPathVariable) SetVariableID(variableID string) {
	u.variableID = variableID
}

func (u URIPathVariable) SetVariableName(variableName string) {
	u.variableName = variableName
}

type PathSegment struct {
	segmentName             string
	templatedSegment        bool
	segmentTemplateDataType SegmentTemplateDataType
}

func (p *PathSegment) GetSegmentName() string {
	return p.segmentName
}

func (p *PathSegment) istemplatedSegment() bool {
	return p.templatedSegment
}

func (p *PathSegment) GetSegmentTemplateDataType() SegmentTemplateDataType {
	return p.segmentTemplateDataType
}

func (p *PathSegment) SetSegmentName(segmentName string) {
	p.segmentName = segmentName
}

func (p *PathSegment) SetTemplatedSegment(templatedSegment bool) {
	p.templatedSegment = templatedSegment
}

func (p *PathSegment) SetSegmentTemplateDataType(segmentTemplateDataType SegmentTemplateDataType) {
	p.segmentTemplateDataType = segmentTemplateDataType
}

func (p PathSegment) equals(object interface{}) bool {
	if p == object {
		return true
	}
	if object == nil {
		return false
	}
	// type check and cast
	otherSegment, ok := object.(PathSegment)

	if !ok {
		return false
	}
	if p.istemplatedSegment() && otherSegment.istemplatedSegment() {
		return p.segmentTemplateDataType == otherSegment.segmentTemplateDataType
	}
	if p.istemplatedSegment() || otherSegment.istemplatedSegment() {
		return false
	}
	return p.segmentName == otherSegment.GetSegmentName()
}

func (p PathSegment) hashCode() int {
	return 79
}

func GetNonTemplatedURI(uriPath string) URI {
	return URI{UriPath: uriPath, HasPathVariable: false}
}

func GetURI(uriPath string) URI {
	pathSegments := GetPathSegments(uriPath)
	isTemplateURI := false
	for _, pathSegment := range pathSegments {
		if isPathSegmentTemplate(pathSegment) {
			isTemplateURI = true
			break
		}
	}
	return URI{UriPath: uriPath, HasPathVariable: isTemplateURI}
}

func (u URI) Equals(object interface{}) bool {
	if u == object {
		return true
	}
	if object == nil {
		return false
	}
	// type check and cast
	otherURI, ok := object.(URI)

	if !ok {
		return false
	}
	if !u.HasPathVariable && !otherURI.HasPathVariable {
		return u.UriPath == otherURI.UriPath
	}

	pathSegments := GetPathSegments(u.UriPath)
	otherURIPathSegments := GetPathSegments(otherURI.UriPath)

	if len(pathSegments) != len(otherURIPathSegments) {
		return false
	}
	for idx := 0; idx < len(pathSegments); idx++ {
		if !arePathSegmentMatching(pathSegments[idx], otherURIPathSegments[idx]) {
			return false
		}
	}
	return true
}

func (u URI) hashCode() int {
	PRIME := 59
	result := 1
	result = result*PRIME + len(GetPathSegments(u.UriPath))
	return result
}

func (u URI) GetSize() int {
	return len(GetPathSegments(u.UriPath))
}

func (u URI) GetPathSegments() []PathSegment {
	pathSegments := GetPathSegments(u.UriPath)
	var pathSegmentList []PathSegment
	for _, pathSegment := range pathSegments {
		pathSegmentList = append(pathSegmentList, GetPathSegment(pathSegment))
	}
	return pathSegmentList
}

func (p PathSegment) String() string {
	return fmt.Sprintf("PathSegment{segmentName=%s, templatedSegment=%t, segmentTemplateDataType=%s}",
		p.segmentName, p.templatedSegment, p.segmentTemplateDataType)
}

func (u URIPathVariable) String() string {
	return fmt.Sprintf("URIPathVariable{variableId=%s, variableName=%s}",
		u.variableID, u.variableName)
}

func (u URI) String() string {
	return fmt.Sprintf("URI{UriPath=%s, HasPathVariable=%t}", u.UriPath, u.HasPathVariable)
}

func (u URI) GetURIPath() string {
	return u.UriPath
}

func (u URI) GetHasPathVariable() bool {
	return u.HasPathVariable
}

func (u *URI) SetURIPath(uriPath string) {
	u.UriPath = uriPath
}

func (u *URI) SetHasPathVariable(hasPathVariable bool) {
	u.HasPathVariable = hasPathVariable
}

//func (u *URI) UnmarshalJSON(data []byte) error {
//
//	// UriPath           string
//	// HasPathVariable   bool
//
//	var tmpJson map[string]interface{}
//
//	if err := json.Unmarshal(data, &tmpJson); err != nil {
//		return err
//	}
//
//	uriPath, ok := tmpJson["UriPath"].(string)
//	if !ok {
//		return errors.New("UriPath is not string")
//	}
//	hasPathVariable, ok := tmpJson["HasPathVariable"].(bool)
//	if !ok {
//		return errors.New("HasPathVariable is not bool")
//	}
//
//	u.SetURIPath(uriPath)
//	u.SetHasPathVariable(hasPathVariable)
//	return nil
//}

func GetPathSegments(uri string) []string {
	if len(uri) == 0 {
		return []string{}
	}
	omitEmptyStrings := func(c rune) bool {
		return c == '/'
	}
	return strings.FieldsFunc(uri, omitEmptyStrings)
}

func GetPathSegment(segmentName string) PathSegment {
	if isPathSegmentTemplate(segmentName) {
		return GetTemplateSegment(GetTemplateDataTypeByDisplayName(segmentName))
	} else {
		return GetNonTemplateSegment(segmentName)
	}
}

func isPathSegmentTemplate(pathSegment string) bool {
	if len(pathSegment) == 0 {
		return false
	}
	trimmedPathSegment := strings.TrimSpace(pathSegment)
	return strings.HasPrefix(trimmedPathSegment, "{") && strings.HasSuffix(trimmedPathSegment, "}")
}

func GetTemplateSegment(segmentTemplateDataType SegmentTemplateDataType) PathSegment {
	return PathSegment{segmentName: segmentTemplateDataType.GetPathDisplayName(),
		templatedSegment:        true,
		segmentTemplateDataType: segmentTemplateDataType}
}

func GetNonTemplateSegment(segmentName string) PathSegment {
	return PathSegment{segmentName: segmentName, templatedSegment: false}
}

func arePathSegmentMatching(pathVariableA string, pathVariableB string) bool {
	if pathVariableA == "" && pathVariableB == "" {
		return true
	}
	if pathVariableA == "" || pathVariableB == "" {
		return false
	}
	if pathVariableA == pathVariableB {
		return true
	}
	isPathVariableATemplate := isPathSegmentTemplate(pathVariableA)
	isPathVariableBTemplate := isPathSegmentTemplate(pathVariableB)
	return isPathVariableATemplate || isPathVariableBTemplate
}
