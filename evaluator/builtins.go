package evaluator

// built-in functions
var bifs = map[string]string{
	// node-test
	"item":                   "item()",
	"node":                   "node()",
	"attribute":              "attribute()",
	"comment":                "comment()",
	"document":               "document()",
	"element":                "element()",
	"schema-element":         "schema-element()",
	"processing-instruction": "processing-instruction()",
	"text":                   "text()",
	"document-node":          "document-node()",

	// func
	"abs":                               "abs()",
	"acos":                              "acos()",
	"add-dayTimeDurations":              "add-dayTimeDurations()",
	"add-dayTimeDuration-to-date":       "add-dayTimeDuration-to-date()",
	"add-dayTimeDuration-to-dateTime":   "add-dayTimeDuration-to-dateTime()",
	"add-dayTimeDuration-to-time":       "add-dayTimeDuration-to-time()",
	"add-yearMonthDurations":            "add-yearMonthDurations()",
	"add-yearMonthDuration-to-date":     "add-yearMonthDuration-to-date()",
	"add-yearMonthDuration-to-dateTime": "add-yearMonthDuration-to-dateTime()",
	"adjust-dateTime-to-timezone":       "adjust-dateTime-to-timezone()",
	"adjust-date-to-timezone":           "adjust-date-to-timezone()",
	"adjust-time-to-timezone":           "adjust-time-to-timezone()",
	"analyze-string":                    "analyze-string()",
	"asin":                              "asin()",
	"atan":                              "atan()",
	"atan2":                             "atan2()",
	"available-environment-variables":   "available-environment-variables()",
	"avg":                               "avg()",
	"base64Binary-equal":                "base64Binary-equal()",
	"base-uri":                          "base-uri()",
	"boolean":                           "boolean()",
	"boolean-equal":                     "boolean-equal()",
	"boolean-greater-than":              "boolean-greater-than()",
	"boolean-less-than":                 "boolean-less-than()",
	"ceiling":                           "ceiling()",
	"codepoint-equal":                   "codepoint-equal()",
	"codepoints-to-string":              "codepoints-to-string()",
	"collection":                        "collection()",
	"compare":                           "compare()",
	"concat":                            "concat()",
	"concatenate":                       "concatenate()",
	"contains":                          "contains()",
	"cos":                               "cos()",
	"count":                             "count()",
	"current-date":                      "current-date()",
	"current-dateTime":                  "current-dateTime()",
	"current-time":                      "current-time()",
	"data":                              "data()",
	"date-equal":                        "date-equal()",
	"date-greater-than":                 "date-greater-than()",
	"date-less-than":                    "date-less-than()",
	"dateTime":                          "dateTime()",
	"dateTime-equal":                    "dateTime-equal()",
	"dateTime-greater-than":             "dateTime-greater-than()",
	"dateTime-less-than":                "dateTime-less-than()",
	"day-from-date":                     "day-from-date()",
	"day-from-dateTime":                 "day-from-dateTime()",
	"days-from-duration":                "days-from-duration()",
	"dayTimeDuration-greater-than":      "dayTimeDuration-greater-than()",
	"dayTimeDuration-less-than":         "dayTimeDuration-less-than()",
	"deep-equal":                        "deep-equal()",
	"default-collation":                 "default-collation()",
	"distinct-values":                   "distinct-values()",
	"divide-dayTimeDuration":            "divide-dayTimeDuration()",
	"divide-dayTimeDuration-by-dayTimeDuration":     "divide-dayTimeDuration-by-dayTimeDuration()",
	"divide-yearMonthDuration":                      "divide-yearMonthDuration()",
	"divide-yearMonthDuration-by-yearMonthDuration": "divide-yearMonthDuration-by-yearMonthDuration()",
	"doc":                                    "doc()",
	"doc-available":                          "doc-available()",
	"document-uri":                           "document-uri()",
	"duration-equal":                         "duration-equal()",
	"element-with-id":                        "element-with-id()",
	"empty":                                  "empty()",
	"encode-for-uri":                         "encode-for-uri()",
	"ends-with":                              "ends-with()",
	"environment-variable":                   "environment-variable()",
	"error":                                  "error()",
	"escape-html-uri":                        "escape-html-uri()",
	"exactly-one":                            "exactly-one()",
	"except":                                 "except()",
	"exists":                                 "exists()",
	"exp":                                    "exp()",
	"exp10":                                  "exp10()",
	"false":                                  "false()",
	"filter":                                 "filter()",
	"floor":                                  "floor()",
	"fold-left":                              "fold-left()",
	"fold-right":                             "fold-right()",
	"for-each":                               "for-each()",
	"for-each-pair":                          "for-each-pair()",
	"format-date":                            "format-date()",
	"format-dateTime":                        "format-dateTime()",
	"format-integer":                         "format-integer()",
	"format-number":                          "format-number()",
	"format-time":                            "format-time()",
	"function-arity":                         "function-arity()",
	"function-lookup":                        "function-lookup()",
	"function-name":                          "function-name()",
	"gDay-equal":                             "gDay-equal()",
	"generate-id":                            "generate-id()",
	"gMonthDay-equal":                        "gMonthDay-equal()",
	"gMonth-equal":                           "gMonth-equal()",
	"gYear-equal":                            "gYear-equal()",
	"gYearMonth-equal":                       "gYearMonth-equal()",
	"has-children":                           "has-children()",
	"head":                                   "head()",
	"hexBinary-equal":                        "hexBinary-equal()",
	"hours-from-dateTime":                    "hours-from-dateTime()",
	"hours-from-duration":                    "hours-from-duration()",
	"hours-from-time":                        "hours-from-time()",
	"id":                                     "id()",
	"idref":                                  "idref()",
	"implicit-timezone":                      "implicit-timezone()",
	"index-of":                               "index-of()",
	"innermost":                              "innermost()",
	"in-scope-prefixes":                      "in-scope-prefixes()",
	"insert-before":                          "insert-before()",
	"intersect":                              "intersect()",
	"iri-to-uri":                             "iri-to-uri()",
	"is-same-node":                           "is-same-node()",
	"lang":                                   "lang()",
	"last":                                   "last()",
	"local-name":                             "local-name()",
	"local-name-from-QName":                  "local-name-from-QName()",
	"log":                                    "log()",
	"log10":                                  "log10()",
	"lower-case":                             "lower-case()",
	"matches":                                "matches()",
	"max":                                    "max()",
	"min":                                    "min()",
	"minutes-from-dateTime":                  "minutes-from-dateTime()",
	"minutes-from-duration":                  "minutes-from-duration()",
	"minutes-from-time":                      "minutes-from-time()",
	"month-from-date":                        "month-from-date()",
	"month-from-dateTime":                    "month-from-dateTime()",
	"months-from-duration":                   "months-from-duration()",
	"multiply-dayTimeDuration":               "multiply-dayTimeDuration()",
	"multiply-yearMonthDuration":             "multiply-yearMonthDuration()",
	"name":                                   "name()",
	"namespace-uri":                          "namespace-uri()",
	"namespace-uri-for-prefix":               "namespace-uri-for-prefix()",
	"namespace-uri-from-QName":               "namespace-uri-from-QName()",
	"nilled":                                 "nilled()",
	"node-after":                             "node-after()",
	"node-before":                            "node-before()",
	"node-name":                              "node-name()",
	"normalize-space":                        "normalize-space()",
	"normalize-unicode":                      "normalize-unicode()",
	"not":                                    "not()",
	"NOTATION-equal":                         "NOTATION-equal()",
	"number":                                 "number()",
	"numeric-add":                            "numeric-add()",
	"numeric-divide":                         "numeric-divide()",
	"numeric-equal":                          "numeric-equal()",
	"numeric-greater-than":                   "numeric-greater-than()",
	"numeric-integer-divide":                 "numeric-integer-divide()",
	"numeric-less-than":                      "numeric-less-than()",
	"numeric-mod":                            "numeric-mod()",
	"numeric-multiply":                       "numeric-multiply()",
	"numeric-subtract":                       "numeric-subtract()",
	"numeric-unary-minus":                    "numeric-unary-minus()",
	"numeric-unary-plus":                     "numeric-unary-plus()",
	"one-or-more":                            "one-or-more()",
	"outermost":                              "outermost()",
	"parse-xml":                              "parse-xml()",
	"parse-xml-fragment":                     "parse-xml-fragment()",
	"path":                                   "path()",
	"pi":                                     "pi()",
	"position":                               "position()",
	"pow":                                    "pow()",
	"prefix-from-QName":                      "prefix-from-QName()",
	"QName":                                  "QName()",
	"QName-equal":                            "QName-equal()",
	"remove":                                 "remove()",
	"replace":                                "replace()",
	"resolve-QName":                          "resolve-QName()",
	"resolve-uri":                            "resolve-uri()",
	"reverse":                                "reverse()",
	"root":                                   "root()",
	"round":                                  "round()",
	"round-half-to-even":                     "round-half-to-even()",
	"seconds-from-dateTime":                  "seconds-from-dateTime()",
	"seconds-from-duration":                  "seconds-from-duration()",
	"seconds-from-time":                      "seconds-from-time()",
	"serialize":                              "serialize()",
	"sin":                                    "sin()",
	"sqrt":                                   "sqrt()",
	"starts-with":                            "starts-with()",
	"static-base-uri":                        "static-base-uri()",
	"string":                                 "string()",
	"string-join":                            "string-join()",
	"string-length":                          "string-length()",
	"string-to-codepoints":                   "string-to-codepoints()",
	"subsequence":                            "subsequence()",
	"substring":                              "substring()",
	"substring-after":                        "substring-after()",
	"substring-before":                       "substring-before()",
	"subtract-dates":                         "subtract-dates()",
	"subtract-dateTimes":                     "subtract-dateTimes()",
	"subtract-dayTimeDuration-from-date":     "subtract-dayTimeDuration-from-date()",
	"subtract-dayTimeDuration-from-dateTime": "subtract-dayTimeDuration-from-dateTime()",
	"subtract-dayTimeDuration-from-time":     "subtract-dayTimeDuration-from-time()",
	"subtract-dayTimeDurations":              "subtract-dayTimeDurations()",
	"subtract-times":                         "subtract-times()",
	"subtract-yearMonthDuration-from-date":   "subtract-yearMonthDuration-from-date()",
	"subtract-yearMonthDuration-from-dateTime": "subtract-yearMonthDuration-from-dateTime()",
	"subtract-yearMonthDurations":              "subtract-yearMonthDurations()",
	"sum":                                      "sum()",
	"tail":                                     "tail()",
	"tan":                                      "tan()",
	"time-equal":                               "time-equal()",
	"time-greater-than":                        "time-greater-than()",
	"time-less-than":                           "time-less-than()",
	"timezone-from-date":                       "timezone-from-date()",
	"timezone-from-dateTime":                   "timezone-from-dateTime()",
	"timezone-from-time":                       "timezone-from-time()",
	"to":                                       "to()",
	"tokenize":                                 "tokenize()",
	"trace":                                    "trace()",
	"translate":                                "translate()",
	"true":                                     "true()",
	"union":                                    "union()",
	"unordered":                                "unordered()",
	"unparsed-text":                            "unparsed-text()",
	"unparsed-text-available":                  "unparsed-text-available()",
	"unparsed-text-lines":                      "unparsed-text-lines()",
	"upper-case":                               "upper-case()",
	"uri-collection":                           "uri-collection()",
	"year-from-date":                           "year-from-date()",
	"year-from-dateTime":                       "year-from-dateTime()",
	"yearMonthDuration-greater-than":           "yearMonthDuration-greater-than()",
	"yearMonthDuration-less-than":              "yearMonthDuration-less-than()",
	"years-from-duration":                      "years-from-duration()",
	"zero-or-one":                              "zero-or-one()",
}

// IsBIF checks if (ident string) is a built-in function or not
func IsBIF(ident string) bool {
	if _, ok := bifs[ident]; ok {
		return true
	}
	return false
}
