/*
Package templatefuncs provides a set of template functions that are used
in the internal template package.
*/
package templatefuncs

func NewFuncs(reporter Reporter) map[string]any {
	return map[string]any{
		"string": get(&StringModule{Reporter: reporter}),
		"lines":  get(&LineModule{}),
		"json":   get(&JSONModule{Reporter: reporter}),
		"yaml":   get(&YAMLModule{Reporter: reporter}),
		"base64": get(&Base64Module{Reporter: reporter}),
		"base32": get(&Base32Module{Reporter: reporter}),
		"gzip":   get(&GZipModule{Reporter: reporter}),
		"list":   get(&ListModule{Reporter: reporter}),

		"html": get(&HTMLModule{}),

		"crc32": get(&CRC32Module{}),
		"crc64": get(&CRC64Module{}),
		"fnv":   get(&FNVModule{}),
		"sha":   get(&SHAModule{}),
	}
}

func get[T any](v T) func() T {
	return func() T { return v }
}
