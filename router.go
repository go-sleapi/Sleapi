package sleapi

import (
	"errors"
	"fmt"
	//"io"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

const (
	LeftBrace          = uint8('{')
	RightBrace         = uint8('}')
	IndexMethod        = "Index"
	StandardParameters = 3
)

type SegmentType int

const (
	PathSeparatorSegmentType = iota
	PathParameterSegmentType = iota
	ContentsSegmentType      = iota
)

type Router struct {
	Routes *RouteTable
}

type Route struct {
	Name         string
	Pattern      string
	Controller   Controller
	RootPattern  string
	Parameter    string
	PathSegments []PathSegmenter
}

type RouteTable struct {
	Table map[string]*Route
}

type RouterDefaults struct {
	Defaults map[string]string
}

type PathSegment struct {
	SegmentContents string
	segmentType     SegmentType
}

type PathSegmenter interface {
	Content() string
	GetSegmentType() SegmentType
}

type SeparatorSegment struct {
	PathSegment
}

type PathParameterSegment struct {
	PathSegment
}

type ContentSegment struct {
	PathSegment
}

type RouterOptions struct {
}

type RouteContext struct {
	Resp       http.ResponseWriter
	Req        *http.Request
	MethodName string
	Param      string
	Data       *RouteData
}

type RouteData struct {
	RouteValues        map[string]string
	PositionParameters []string
}

func NewRouteTable() *RouteTable {
	rt := &RouteTable{}
	rt.Table = make(map[string]*Route)
	return rt
}

func NewRouter() *Router {
	rt := NewRouteTable()
	return &Router{rt}
}

func NewRouteContext(w http.ResponseWriter, req *http.Request) *RouteContext {
	rc := &RouteContext{w, req, "", "", nil}
	return rc
}

func NewRoute(name string, routePattern string, controller Controller, root string, param string) *Route {
	r := &Route{name, routePattern, controller, root, param, nil}
	r.PathSegments = make([]PathSegmenter, 0)
	return r
}

func NewSeparatorSegment() *SeparatorSegment {
	ss := &SeparatorSegment{}
	ss.SegmentContents = "/"
	ss.segmentType = PathSeparatorSegmentType
	return ss
}

func NewPathParameterSegment(contents string) *PathParameterSegment {
	//fmt.Println("Creating Path Parameter Segment...", contents)
	pps := &PathParameterSegment{}
	pps.SegmentContents = contents
	pps.segmentType = PathParameterSegmentType
	return pps
}

func NewContentSegment(contents string) *ContentSegment {
	//fmt.Println("Creating Content Segment...", contents)
	cs := &ContentSegment{}
	cs.SegmentContents = contents
	cs.segmentType = ContentsSegmentType
	return cs
}

func NewRouteData() *RouteData {
	rd := &RouteData{}
	rd.RouteValues = make(map[string]string)
	rd.PositionParameters = make([]string, 0)
	return rd
}

func (this PathSegment) Content() string {
	return this.SegmentContents
}

func (this PathSegment) GetSegmentType() SegmentType {
	return this.segmentType
}

func (this PathSegment) String() string {
	//return this.Content()
	return strconv.Itoa(int(this.segmentType))
}

func (this PathParameterSegment) String() string {
	//return "{" + this.SegmentContents + "}"
	return strconv.Itoa(int(this.segmentType))
}

func (this PathParameterSegment) GetSegmentType() SegmentType {
	return PathParameterSegmentType
}

/*func (this SeparatorSegment) Content() string {
	return this.SegmentContents
}
*/

func parseParameter(pattern string) (string, string) {
	lb := strings.Index(pattern, string(LeftBrace))
	/*if lb != -1 {
		fmt.Println("Found { at: ", lb)
	}
	*/

	var token string
	var root string
	for i := lb + 1; i < len(pattern); i++ {
		if pattern[i] == RightBrace {
			break
		}
		token += string(pattern[i])
	}

	sandbox := pattern
	if sandbox[0] == '/' {
		sandbox = sandbox[1:]
	}
	parts := strings.Split(sandbox, "/")

	//fmt.Println("Parts(parseParameter): ", parts)

	for _, part := range parts {
		if !strings.Contains(part, "{") && !strings.Contains(part, "}") {
			root += "/" + part
		}
	}

	return token, root
}

func parseRoute(requestPath string) []PathSegmenter {
	parts := splitIntoPathStrings(requestPath)
	//fmt.Println("Parts: ", parts)

	segments := splitIntoPathSegments(parts)
	//fmt.Println("Segments: ", segments)

	return segments
}

func (this *Router) AddRoute(name string, routePattern string, controller Controller) {
	param, root := parseParameter(routePattern)
	segments := parseRoute(routePattern)
	route := NewRoute(name, routePattern, controller, root, param)
	route.PathSegments = segments

	this.Routes.Table[name] = route
}

func splitIntoPathStrings(requestPath string) []string {
	//fmt.Println("requestPath", requestPath)
	parts := make([]string, 0)

	if len(requestPath) == 0 {
		return parts
	}

	currentIndex := 0

	for currentIndex < len(requestPath) {
		curPath := requestPath[currentIndex:]
		//fmt.Println("curPath: ", curPath)
		indexOfNextSeparator := strings.Index(curPath, "/")
		if indexOfNextSeparator == -1 {
			finalPart := requestPath[currentIndex:]
			if len(finalPart) > 0 {
				parts = append(parts, finalPart)
			}
			break
		}

		//fmt.Println("currentIndex: ", currentIndex)
		//fmt.Println("indexOfNextSeparator: ", indexOfNextSeparator)
		nextPart := curPath[0:indexOfNextSeparator]
		//fmt.Println("nextPart: ", nextPart)
		if len(nextPart) > 0 {
			parts = append(parts, nextPart)
		}
		parts = append(parts, "/")
		currentIndex = currentIndex + indexOfNextSeparator + 1
	}

	return parts
}

func splitIntoPathSegments(parts []string) []PathSegmenter {
	segments := make([]PathSegmenter, 0)

	for _, part := range parts {
		isSeparator := isSeparator(part)
		if isSeparator {
			sp := NewSeparatorSegment()
			segments = append(segments, sp)
		} else {
			segment := parseSegment(part)
			//fmt.Println("Segment Type: ", segment.GetSegmentType())
			segments = append(segments, segment)
			//fmt.Println("Segments (split): ", segments)
		}
	}

	return segments
}

func isSeparator(part string) bool {
	return part == "/"
}

func parseSegment(segment string) PathSegmenter {

	startIndex := 0

	//for startIndex < len(part) {
	nextParameterStart := indexOfFirstOpenParameter(segment, startIndex)
	if nextParameterStart == -1 {
		literal := getLiteral(segment[startIndex:])
		if len(literal) > 0 {
			cs := NewContentSegment(literal)
			return cs
		}
	}

	nextParameterEnd := IndexOf(segment, "}", nextParameterStart+1)
	if nextParameterEnd != -1 {
		literal := getLiteral(segment[startIndex : nextParameterStart-startIndex])
		if len(literal) > 0 {
			cs := NewContentSegment(literal)
			return cs
		}

		//parameterName := segment[nextParameterStart + 1:nextParameterEnd - nextParameterStart + 1]
		parameterName := segment[nextParameterStart+1 : nextParameterEnd+1]
		if len(parameterName) > 0 {
			ps := NewPathParameterSegment(parameterName)
			return ps
		}
	}
	//}
	return PathSegment{}
}

func indexOfFirstOpenParameter(segment string, startIndex int) int {
	index := 0
	for {
		index = IndexOf(segment, "{", startIndex)
		if index == -1 {
			return -1
		}

		if (index+1 == len(segment)) || ((index+1 < len(segment)) && (segment[index+1] != '{')) {
			return index
		}

		index += 2
	}
}

func IndexOf(source string, sep string, startIndex int) int {
	thisSlice := source[startIndex:]
	return strings.Index(thisSlice, sep)
}

func getLiteral(segment string) string {
	temp := strings.Replace(segment, "{{", "", -1)
	newLiteral := strings.Replace(temp, "}}", "", -1)
	if strings.Contains(newLiteral, "{{") || strings.Contains(newLiteral, "}}") {
		return ""
	}

	temp = strings.Replace(newLiteral, "{{", "{", -1)
	newLiteral = strings.Replace(temp, "}}", "}", -1)
	return newLiteral
}

func callControllerMethod(route *Route, rc *RouteContext) {
	typ := reflect.ValueOf(route.Controller)
	//fmt.Println("Type: ", typ)

	/*if len(rc.MethodName) > 0 {
		fmt.Println("CallMethod: ", rc.MethodName)
	}
	*/

	me := typ.MethodByName(rc.MethodName)
	if me.IsValid() {
		//fmt.Println("IsValid...")
		in := []reflect.Value{reflect.ValueOf(rc.Resp), reflect.ValueOf(rc.Req)}
		if len(rc.Data.PositionParameters) > 0 {
			for _, p := range rc.Data.PositionParameters {
				//fmt.Println("Appending Val: ", p)
				in = append(in, reflect.ValueOf(p))
			}
		}
		log.Printf("Before calling: %s.%s\n", route.Controller.Name(), rc.MethodName)
		me.Call(in)
		log.Printf("After Calling: %s.%s\n", route.Controller.Name(), rc.MethodName)
	}
}

func isIndexRoute(route *Route, rc *RouteContext) bool {

	pattern := route.RootPattern
	cleanPath := strings.TrimSuffix(rc.Req.URL.Path, "/")
	if cleanPath == pattern && rc.Req.Method == "GET" {
		rc.MethodName = IndexMethod
		//fmt.Println("Found Index Route")
		return true
	}

	return false
}

func findControllerMethod(route *Route, rc *RouteContext) (bool, string, error) {
	typ := reflect.TypeOf(route.Controller)
	//fmt.Println("Type: ", typ)
	numMethods := typ.NumMethod()
	foundMethod := ""
	for i := 0; i < numMethods; i++ {
		method := typ.Method(i)
		//fmt.Println("Method Name: ", method.Name)
		methodName := strings.ToLower(method.Name)
		httpMethod := strings.ToLower(rc.Req.Method)
		numParameters := method.Type.NumIn()
		//fmt.Println("In: ", numParameters)

		//paramPos := 0
		if numParameters == (len(rc.Data.RouteValues)+StandardParameters) && strings.Contains(methodName, httpMethod) {
			//fmt.Println("Method: ", methodName)
			if len(foundMethod) > 0 {
				return false, "", errors.New("Found Multiple Matching Methods")
			}
			foundMethod = method.Name
		}
		/*for p := 0; p < numParameters; p++ {
			param := method.Type.In(p)
			fmt.Println("In Param: ", param)
			fmt.Println("In Param Name: ", param.Name())
			fmt.Println("In Param Name (String()): ", param.String())
			if param.String() != "http.ResponseWriter" && param.String() != "*http.Request" {
				switch param.String() {
				case "string":
					p = rc.Data.PositionParameters[paramPos]

				}
			}
		}
		if strings.Contains(methodName, httpMethod) {
			//rc.MethodName = method.Name
			//rc.Param = param
			return true, method.Name
		}
		*/
	}

	if len(foundMethod) > 0 {
		return true, foundMethod, nil
	}

	return false, "", nil
}

func isMatch(route *Route, rc *RouteContext) bool {
	path := rc.Req.URL.Path
	lastSlashIndex := strings.LastIndex(path, "/")
	if lastSlashIndex != -1 {
		param := path[lastSlashIndex+1:]
		//fmt.Println("Parameter: ", param)
		typ := reflect.TypeOf(route.Controller)
		//fmt.Println("Type: ", typ)
		numMethods := typ.NumMethod()
		for i := 0; i < numMethods; i++ {
			method := typ.Method(i)
			//fmt.Println("Method Name: ", method.Name)
			methodName := strings.ToLower(method.Name)
			httpMethod := strings.ToLower(rc.Req.Method)
			if strings.Contains(methodName, httpMethod) {
				rc.MethodName = method.Name
				rc.Param = param
				return true
			}
		}
	}

	return false
}

func matchRoute(route *Route, segments []string) *RouteData {
	pathSegments := route.PathSegments

	//fmt.Println("matchRoute.Segments: ", segments)

	if len(route.PathSegments) != len(segments) {
		return nil
	}

	routeData := NewRouteData()
	for index, segment := range pathSegments {
		//fmt.Println("Parameter Check...", segment.GetSegmentType())
		//fmt.Println("Parameter Type: ", PathParameterSegmentType)
		if segment.GetSegmentType() != PathParameterSegmentType {
			//fmt.Println("segment.Content(): ", segment.Content())
			//fmt.Println("segments[index]: ", segments[index])
			if segment.Content() != segments[index] {
				return nil
			}
		} else {
			//fmt.Println("{Parameter} = ", segment.Content())
			//fmt.Println("Parameter Value = ", segments[index])
			/*if segment.Content() != segments[index] {
				return nil
			}
			*/
			routeData.RouteValues[segment.Content()] = segments[index]
			routeData.PositionParameters = append(routeData.PositionParameters, segments[index])
		}
	}

	return routeData
}

func (this *Router) FindRoute(w ResponseWriter, req *http.Request) {
	//fmt.Println("Find Route: ", this.Routes.Table)

	//cleanPath := strings.TrimPrefix(req.URL.Path, "/")
	cleanPath := strings.TrimPrefix(req.RequestURI, "/")
	segments := splitIntoPathStrings(cleanPath)
	//fmt.Println("Request Parts: ", segments)

	ok := false
	rc := NewRouteContext(w, req)
	for _, route := range this.Routes.Table {
		fmt.Println("Route Name: ", name)
		//pattern := route.RootPattern

		/*fmt.Println("Root Pattern: ", pattern)
		if isIndexRoute(route, rc) {
			callControllerMethod(route, rc)
			return route
		} else if isMatch(route, rc) {
			fmt.Println("Found Match...")
			callControllerMethod(route, rc)
		}
		*/

		routeData := matchRoute(route, segments)
		if routeData != nil {
			//fmt.Println("Found Match With Segments...", routeData.RouteValues)
			rc.Data = routeData
			found, methodName, err := findControllerMethod(route, rc)
			if err != nil {
				fmt.Println("Error: ", err.Error())
			} else if found {
				//fmt.Println("Found Method: ", methodName)
				rc.MethodName = methodName
				callControllerMethod(route, rc)
				ok = routeData != nil
				return
			}
		}
	}

	if !ok {
		fmt.Println("Route Not Found: 404")
		//w.WriteHeader(http.StatusNotFound)
		//io.WriteString(w, "404 - Page Not Found!\n")
		http.NotFound(w, req)
	}
	//return nil
}
