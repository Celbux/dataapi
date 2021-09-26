// Package i contains a set of commonly used interfaces .The package name is
// kept short so that imports look standard eg. i.Logger still looks like
// iLogger. Try to NOT add interfaces to this package, we only want fundamental
// interfaces here
package i

// Logger is an interface that allows us to decouple
// from the default logging package throughout our code.
// This allows for easier testing and control of data flow.
type Logger interface {
	Println(v ...interface{})
	Printf(format string, v ...interface{})
}
