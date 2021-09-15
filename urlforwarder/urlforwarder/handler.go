package urlforwarder

import (
	"gopkg.in/yaml.v2"
	"net/http"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// // http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		redirectUrl, ok := pathsToUrls[request.URL.Path]
		if ok {
			http.Redirect(writer, request, redirectUrl, http.StatusMovedPermanently)
		} else {
			fallback.ServeHTTP(writer, request)
		}
	}
}

type PathUrl struct {
	Path string
	Url  string
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yamlData []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathUrl, err := parseYAML(yamlData)
	if err != nil {
		return nil, err
	}

	pathMap := buildMap(pathUrl)
	return MapHandler(pathMap, fallback), nil
}

func parseYAML(yamlData []byte) ([]PathUrl, error) {
	var p []PathUrl
	if err := yaml.Unmarshal(yamlData, &p); err != nil {
		return nil, err
	}

	return p, nil
}

func buildMap(pathUrl []PathUrl) map[string]string {
	m := make(map[string]string)
	if pathUrl == nil {
		return m
	}

	for _, pu := range pathUrl {
		_, ok := m[pu.Path]
		if !ok {
			m[pu.Path] = pu.Url
		}
	}

	return m
}
