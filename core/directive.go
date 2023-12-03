package core

import (
	"errors"
	"fmt"

	"github.com/ggicci/httpin/internal"
	"github.com/ggicci/owl"
)

func init() {
	RegisterDirective("form", &DirectvieForm{})
	RegisterDirective("query", &DirectiveQuery{})
	RegisterDirective("header", &DirectiveHeader{})
	RegisterDirective("body", &DirectiveBody{})
	RegisterDirective("required", &DirectiveRequired{})
	RegisterDirective("default", &DirectiveDefault{})
}

// RegisterFileType registers a FileEncodeDecoder for type T. Which marks the type T as
// a file type. When httpin encounters a field of type T, it will treat it as a file
// upload.
//
//	func init() {
//	    RegisterFileType[MyFile](&myFileEncodeDecoder{})
//	}
func RegisterFileType[T FileEncoder](fd FileDecoder[T]) {
	internal.PanicOnError(
		defaultRegistry.RegisterFileType(
			internal.TypeOf[T](),
			toAnyFileDecoder[T](fd),
		),
	)
}

var (
	// decoderNamespace is the namespace for registering directive executors that are
	// used to decode the http request to input struct.
	decoderNamespace = owl.NewNamespace()

	// encoderNamespace is the namespace for registering directive executors that are
	// used to encode the input struct to http request.
	encoderNamespace = owl.NewNamespace()

	reservedExecutorNames = []string{"decoder", "encoder"}
	noopDirective         = &directiveNoop{}
)

type DirectiveExecutor interface {
	// Encode encodes the field of the input struct to the HTTP request.
	Encode(*DirectiveRuntime) error

	// Decode decodes the field of the input struct from the HTTP request.
	Decode(*DirectiveRuntime) error
}

func init() {
	// decoder is a special executor which does nothing, but is an indicator of
	// overriding the decoder for a specific field.
	decoderNamespace.RegisterDirectiveExecutor("decoder", asOwlDirectiveExecutor(noopDirective.Decode))
	encoderNamespace.RegisterDirectiveExecutor("encoder", asOwlDirectiveExecutor(noopDirective.Encode))
}

// RegisterDirective registers a DirectiveExecutor with the given directive name. The
// directive should be able to both extract the value from the HTTP request and build
// the HTTP request from the value. The Decode API is used to decode data from the HTTP
// request to a field of the input struct, and Encode API is used to encode the field of
// the input struct to the HTTP request.
//
// Will panic if the name were taken or given executor is nil. Pass parameter force
// (true) to ignore the name conflict.
func RegisterDirective(name string, executor DirectiveExecutor, force ...bool) {
	registerDirectiveExecutorToNamespace(decoderNamespace, name, executor, force...)
	registerDirectiveExecutorToNamespace(encoderNamespace, name, executor, force...)
}

func registerDirectiveExecutorToNamespace(ns *owl.Namespace, name string, exe DirectiveExecutor, force ...bool) {
	panicOnReservedExecutorName(name)
	if exe == nil {
		internal.PanicOnError(errors.New("nil directive executor"))
	}
	if ns == decoderNamespace {
		ns.RegisterDirectiveExecutor(name, asOwlDirectiveExecutor(exe.Decode), force...)
	} else {
		ns.RegisterDirectiveExecutor(name, asOwlDirectiveExecutor(exe.Encode), force...)
	}
}

func asOwlDirectiveExecutor(directiveFunc func(*DirectiveRuntime) error) owl.DirectiveExecutor {
	return owl.DirectiveExecutorFunc(func(dr *owl.DirectiveRuntime) error {
		return directiveFunc((*DirectiveRuntime)(dr))
	})
}

func panicOnReservedExecutorName(name string) {
	for _, reservedName := range reservedExecutorNames {
		if name == reservedName {
			internal.PanicOnError(fmt.Errorf("reserved executor name: %q", name))
		}
	}
}

// directiveNoop is a DirectiveExecutor that does nothing, "noop" stands for "no operation".
type directiveNoop struct{}

func (*directiveNoop) Encode(*DirectiveRuntime) error { return nil }
func (*directiveNoop) Decode(*DirectiveRuntime) error { return nil }
