package file

import (
	"reflect"
	"strings"
	"testing"

	"github.com/opencharly/plugin-file/candy/plugin-file/params"
)

// TestFileCheckProjectsEveryInputField guards a silent-failure class, not a hypothetical.
//
// `fileCheck` is a HAND-WRITTEN projection of params.FileInput, built by a struct literal
// in RunVerb. A field added to the CUE schema flows automatically into params.FileInput
// when the generator runs — and then reaches NOTHING, because the literal does not copy it.
// Everything still compiles. `charly box validate` still accepts the authored step. The
// check simply ignores the field and passes, which is the worst available outcome: an
// assertion that reads as if it is testing something and is not.
//
// This happened while adding `link_target`: the schema, the generated params type and the
// probe were all correct, and the field was still invisible until the literal was updated.
//
// So: every exported field of params.FileInput must have a same-named field on fileCheck.
// The one legitimate exception is the primary `File`, which the projection renames to
// `Path`.
func TestFileCheckProjectsEveryInputField(t *testing.T) {
	inT := reflect.TypeOf(params.FileInput{})
	chkT := reflect.TypeOf(fileCheck{})

	// Fields the projection deliberately renames or does not carry, each with the reason.
	renamed := map[string]string{
		"File": "Path", // the verb's primary field is spelled Path on the check struct
	}

	if inT.NumField() < 5 {
		t.Fatalf("params.FileInput has only %d fields — the reflection is not seeing the "+
			"real type, so this guard would pass vacuously", inT.NumField())
	}

	for i := 0; i < inT.NumField(); i++ {
		name := inT.Field(i).Name
		if !inT.Field(i).IsExported() {
			continue
		}
		want := name
		if alias, ok := renamed[name]; ok {
			want = alias
		}
		if _, found := chkT.FieldByName(want); !found {
			t.Errorf("params.FileInput.%s has no %s on fileCheck: the field is decoded "+
				"from the step and then DROPPED, so an authored `%s:` silently asserts "+
				"nothing", name, want, strings.ToLower(name))
		}
	}
}

// The mirror: a field on fileCheck that no input field feeds is dead weight the probe may
// be branching on, which reads as a working assertion driven by a permanently-zero value.
func TestFileCheckHasNoFieldTheInputCannotFill(t *testing.T) {
	inT := reflect.TypeOf(params.FileInput{})
	chkT := reflect.TypeOf(fileCheck{})

	fromInput := map[string]bool{"Path": true} // Path is File, renamed
	for i := 0; i < inT.NumField(); i++ {
		fromInput[inT.Field(i).Name] = true
	}
	for i := 0; i < chkT.NumField(); i++ {
		name := chkT.Field(i).Name
		if !fromInput[name] {
			t.Errorf("fileCheck.%s is fed by no params.FileInput field — the probe branches "+
				"on a value nothing can set", name)
		}
	}
}
