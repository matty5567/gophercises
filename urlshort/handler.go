package urlshort

import (
	"net/http"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	redirectFunc := func(w http.ResponseWriter, r *http.Request) {
		for short_url, long_url := range pathsToUrls {
			if r.URL.Path == short_url {
				http.Redirect(w, r, long_url, http.StatusSeeOther)
				return
			}
		}
		fallback.ServeHTTP(w, r)

	}

	return http.HandlerFunc(redirectFunc)
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {

	type Config struct {
		Path string
		Url  string
	}

	var configs []Config
	err := yaml.Unmarshal(yml, &configs)

	if err != nil {
		return nil, err
	}

	redirect_map := make(map[string]string)

	for _, config := range configs {
		redirect_map[config.Path] = config.Url
	}

	return MapHandler(redirect_map, fallback), nil
}
