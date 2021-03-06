package parser

import (
	"reflect"
	"testing"

	"lang/pkg/parser"

	"github.com/alecthomas/participle"
	"github.com/stretchr/testify/assert"
)

type TestVal struct {
	Field  string
	ValStr string
	ValInt int32
}

// Only useful for testing function declarations with a string as the body.
type TestFnDeclVal struct {
	Name string
	Args []string
	Body string
}

type TestFnCallVal struct {
	Name string
	Args []TestVal
}

type TestVarDeclVal struct {
	Name  string
	Value string
}

func TestMainPrimitives(t *testing.T) {
	// Build the assertor and the tokenizer
	tokenizer, err := participle.Build(&parser.Program{})
	assert := assert.New(t)

	if err != nil {
		panic(err)
	}

	// Holds in the values for later
	// Each program is in the form
	// @program:<@field,@value>
	programs := map[string]TestVal{
		`(main () true)`: {
			Field:  "Atom",
			ValStr: "true",
		},
		`(main () b101101)`: {
			Field:  "Atom",
			ValStr: "b101101",
		},
		`(main () 1)`: {
			Field:  "Int",
			ValInt: 1,
		},
		`(main () 3)`: {
			Field:  "Int",
			ValInt: 3,
		},
		`(main () 1000033)`: {
			Field:  "Int",
			ValInt: 1000033,
		},
		`(main () "1")`: {
			Field:  "Str",
			ValStr: "1",
		},
		`(main () "JOE_IS_CoOL!!1")`: {
			Field:  "Str",
			ValStr: "JOE_IS_CoOL!!1",
		},
	}

	for sourceCode, output := range programs {
		fieldName := output.Field

		root := &parser.Program{}
		tokenizer.ParseString(sourceCode, root)

		// Test the non nillness of main.
		assert.NotNilf(root.DefOrMain[0].Main, "Hmm, main is nil for `%s`", sourceCode)
		// Test if parsed value is expected.
		parsedOutput := reflect.ValueOf(*root.DefOrMain[0].Main.Body).FieldByName(fieldName).String()

		// We must treat ints different to strings
		if parsedOutput == "<int32 Value>" {
			intOutput := reflect.ValueOf(*root.DefOrMain[0].Main.Body).FieldByName(fieldName).Int()
			assert.Equalf(
				output.ValInt,
				int32(intOutput), // This conversion is because reflect outputs to 64 bit ints
				"Hmm, we failed the int output, for: %s", sourceCode,
			)
		} else {
			assert.Equalf(
				output.ValStr, parsedOutput,
				"Hmm, we failed the string output, for: %s", sourceCode,
			)
		}
	}
}

func TestFnDecl(t *testing.T) {
	// Build the assertor and the tokenizer
	tokenizer, err := participle.Build(&parser.Program{})
	assert := assert.New(t)

	if err != nil {
		panic(err)
	}

	// Holds in the values for later
	// Each program is in the form
	// @program:<@Name,@Args,@Body>
	programs := map[string]TestFnDeclVal{
		`(defun test (a b) "abc")`: {
			Name: "test",
			Args: []string{"a", "b"},
			Body: "abc",
		},
	}

	for sourceCode, output := range programs {
		root := &parser.Program{}
		tokenizer.ParseString(sourceCode, root)
		// Test the non nilness of FnDecl.
		assert.NotNilf(root.DefOrMain[0].FnDecl, "Hmm, FnDecl is nil for `%s`", sourceCode)
		// Test function name
		parsedOutput := reflect.ValueOf(*root.DefOrMain[0].FnDecl).FieldByName("Name").String()
		assert.Equalf(output.Name, parsedOutput, "Hmm, we failed the function declaration name, for: %s", sourceCode)
		// Test length of args.
		assert.Equal(len(root.DefOrMain[0].FnDecl.Args), len(output.Args), "Hmm, the number of args for FnDecl is incorrect.`%s`", sourceCode)
		// Test Args
		for i, value := range output.Args {
			parsedOutput = reflect.ValueOf(*root.DefOrMain[0].FnDecl.Args[i]).FieldByName("Atom").String()
			assert.Equalf(value, parsedOutput, "Hmm, we failed the function declaration args, for: %s", sourceCode)
		}
		// Test Body
		parsedOutput = reflect.ValueOf(*root.DefOrMain[0].FnDecl.Body).FieldByName("Str").String()
		assert.Equalf(output.Body, parsedOutput, "Hmm, we failed the function declaration body, for: %s", sourceCode)
	}
}

func TestFnCall(t *testing.T) {
	// Build the assertor and the tokenizer
	tokenizer, err := participle.Build(&parser.Program{})
	assert := assert.New(t)

	if err != nil {
		panic(err)
	}

	// Holds in the values for later
	// Each program is in the form
	// @program:<@Name,@Args>
	programs := map[string]TestFnCallVal{
		`(main () (eat 1 2 3))`: {
			Name: "eat",
			Args: []TestVal{
				{
					Field:  "int",
					ValInt: 1,
				},
				{
					Field:  "int",
					ValInt: 2,
				},
				{
					Field:  "int",
					ValInt: 3,
				},
			},
		},
	}

	for sourceCode, output := range programs {
		root := &parser.Program{}
		tokenizer.ParseString(sourceCode, root)
		// Test the non nilness of FnCall.
		assert.NotNilf(root.DefOrMain[0].Main.Body.Fn, "Hmm, FnCall is nil for `%s`", sourceCode)
		// Test function name
		parsedOutput := reflect.ValueOf(*root.DefOrMain[0].Main.Body.Fn).FieldByName("Name").String()
		assert.Equalf(output.Name, parsedOutput, "Hmm, we failed the function call name, for: %s", sourceCode)
		// Test length of args.
		assert.Equal(len(root.DefOrMain[0].Main.Body.Fn.Args), len(output.Args), "Hmm, the number of args for FnCall is incorrect.`%s`", sourceCode)
		// Test Args
		for i, value := range output.Args {
			// Test if parsed value is expected.
			parsedOutput = reflect.ValueOf(*root.DefOrMain[0].Main.Body.Fn.Args[i]).FieldByName("Str").String()

			// We must treat ints different to strings
			if parsedOutput == "<int32 Value>" {
				intOutput := reflect.ValueOf(*root.DefOrMain[0].Main.Body.Fn.Args[i]).FieldByName("Int").Int()
				assert.Equalf(
					value.ValInt, // This conversion is because reflect outputs to 64 bit ints
					int32(intOutput),
					"Hmm, we failed the function call args (int), for: %s", sourceCode,
				)
			} else {
				assert.Equalf(
					value.ValStr, parsedOutput,
					"Hmm, we failed the function call args (string), for: %s", sourceCode,
				)
			}
		}
	}
}

func TestVarDecl(t *testing.T) {
	// Build the assertor and the tokenizer
	tokenizer, err := participle.Build(&parser.Program{})
	assert := assert.New(t)

	if err != nil {
		panic(err)
	}

	// Holds in the values for later
	// Each program is in the form
	// @program:<@Name,@Value>
	programs := map[string]TestVarDeclVal{
		`(defvar ALICE_ADDR "1LCZTUkMKSYN8oKWhh8oqTErEhTENpnXY6")`: {
			Name:  "ALICE_ADDR",
			Value: "1LCZTUkMKSYN8oKWhh8oqTErEhTENpnXY6",
		},
	}

	for sourceCode, output := range programs {
		root := &parser.Program{}
		tokenizer.ParseString(sourceCode, root)
		// Test the non nilness of VarDecl.
		assert.NotNilf(root.DefOrMain[0].VarDecl, "Hmm, VarDecl is nil for `%s`", sourceCode)
		// Test variable name
		parsedOutput := reflect.ValueOf(*root.DefOrMain[0].VarDecl).FieldByName("Name").String()
		assert.Equalf(output.Name, parsedOutput, "Hmm, we failed the variable declaration name, for: %s", sourceCode)
		// Test variable value
		parsedOutput = reflect.ValueOf(*root.DefOrMain[0].VarDecl.Value).FieldByName("Str").String()
		assert.Equalf(output.Value, parsedOutput, "Hmm, we failed the variable declaration value, for: %s", sourceCode)
	}
}
