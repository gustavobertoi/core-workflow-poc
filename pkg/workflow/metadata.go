package workflow

type NodeMetadata interface {
	Input() InputOutputMetadata
	Output() InputOutputMetadata
}

type nodeMetadata struct {
	input  InputOutputMetadata
	output InputOutputMetadata
}

func NewNodeMetadata(input InputOutputMetadata, output InputOutputMetadata) NodeMetadata {
	return &nodeMetadata{
		input:  input,
		output: output,
	}
}

func (n *nodeMetadata) Input() InputOutputMetadata {
	return n.input
}

func (n *nodeMetadata) Output() InputOutputMetadata {
	return n.output
}

type InputOutputMetadata struct {
	Parameters Parameters
	Edges      EdgeMetadata
}

type EdgeMetadata struct {
	RouteName  string
	Count      EdgeCount
	Parameters Parameters
}

type EdgeCount int

const (
	EdgesUnlimited EdgeCount = -1
)

type Parameters map[string]ParameterSchema

// ParameterSchema represents a schema definition for a single data field.
// Each field in the schema can have specific properties like type, validation rules, and metadata.
//
// Validations array format:
// The Validations slice contains strings that specify rules for validation.
// Examples:
//
// - "min=18": Ensures the field has a minimum value of 18 (applicable to numeric types like int or float).
// - "max=99": Ensures the field has a maximum value of 99 (applicable to numeric types like int or float).
// - "regex=^[a-zA-Z0-9]+$": Ensures the field matches a specific regular expression pattern (applicable to string types).
// - "len=10": Ensures the field is exactly 10 characters long (for strings or arrays).
// - "in=male,female,other": Ensures the field value is one of the specified allowed values.
// - "required": Ensures the field is mandatory, though typically this is also expressed with the Required bool field.
// - "email": Ensures the field contains a valid email address (applicable to string types).
// - "uuid": Ensures the field contains a valid UUID value.
//
// Example usage:
// FieldName: "Username", Type: "string", Required: true, Validations: []string{"len=8", "regex=^[a-zA-Z]+$"}
// FieldName: "Age", Type: "int", Required: true, Validations: []string{"min=18", "max=65"}
type ParameterSchema struct {
	Name        string   // The variable name
	Type        string   // The variable type (e.g., string, int, bool, etc.)
	Required    bool     // Whether the field is mandatory
	Validations []string // A list of validations to apply (e.g., min, max, regex, etc.)
	Description string   // Optional description of the field
	Default     any      // Default value if any
}
