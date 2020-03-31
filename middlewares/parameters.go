package middlewares

import (
	"bytes"
	"net/http"
	"regexp"
	"strings"
)

func (p pattern) WithParameters(parameters ...string) pattern {
	decorate := parseParameters(parameters, string(p))
	patternWithoutParams := string(p)
	ro := routes[string(p)]
	delete(routes, string(p))
	decorated := decorate(ro)
	for _, param := range parameters {
		replace := "/{" + param + "}"
		patternWithoutParams = strings.Replace(patternWithoutParams, replace, "", -1)
	}
	routes[patternWithoutParams+"/"] = decorated
	return pattern(patternWithoutParams + "/")
}

func parseParameters(parameters []string, p string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//send the parameter in request
			patternWithoutParams := p
			for _, param := range parameters {
				b := bytes.Buffer{}
				b.WriteString("{")
				b.WriteString(param)
				b.WriteString("}")
				replace := b.String()
				patternWithoutParams = strings.Replace(patternWithoutParams, replace, "([^/]*)", -1)
			}

			reg, err := regexp.Compile(patternWithoutParams)

			if err != nil {
				return
			}

			stringToMatch := r.URL.Path

			s := strings.Split(stringToMatch, "/")

			paramLength := len(s) - 2

			if paramLength < len(parameters) {
				stringToMatch = stringToMatch + "/"
			}

			matches := reg.FindAllStringSubmatch(stringToMatch, -1)

			if matches != nil {
				matchSlice := matches[0]
				values := r.URL.Query()
				for i, _ := range matchSlice {
					if i > 0 && matchSlice[i] != "" {
						values.Add(parameters[i-1], matchSlice[i])
					}
				}
				r.URL.RawQuery = values.Encode()
			}

			h.ServeHTTP(w, r)
		})
	}
}
