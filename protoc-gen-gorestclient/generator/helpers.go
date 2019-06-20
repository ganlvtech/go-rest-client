package generator

import (
	"log"
	"strings"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

// Is c an ASCII lower-case letter?
func isASCIILower(c byte) bool {
	return 'a' <= c && c <= 'z'
}

// Is c an ASCII digit?
func isASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

// CamelCase returns the CamelCased name.
// If there is an interior underscore followed by a lower case letter,
// drop the underscore and convert the letter to upper case.
// There is a remote possibility of this rewrite causing a name collision,
// but it's so remote we're prepared to pretend it's nonexistent - since the
// C++ generator lowercases names, it's extremely unlikely to have two fields
// with different capitalizations.
// In short, _my_field_name_2 becomes XMyFieldName_2.
func CamelCase(s string) string {
	if s == "" {
		return ""
	}
	t := make([]byte, 0, 32)
	i := 0
	if s[0] == '_' {
		// Need a capital letter; drop the '_'.
		t = append(t, 'X')
		i++
	}
	// Invariant: if the next letter is lower case, it must be converted
	// to upper case.
	// That is, we process a word at a time, where words are marked by _ or
	// upper case letter. Digits are treated as words.
	for ; i < len(s); i++ {
		c := s[i]
		if c == '_' && i+1 < len(s) && isASCIILower(s[i+1]) {
			continue // Skip the underscore in s.
		}
		if isASCIIDigit(c) {
			t = append(t, c)
			continue
		}
		// Assume we have a letter now - if not, it's a bogus identifier.
		// The next word is a sequence of characters that must start upper case.
		if isASCIILower(c) {
			c ^= ' ' // Make it a capital letter.
		}
		t = append(t, c) // Guaranteed not lower case.
		// Accept lower case sequence that follows.
		for i+1 < len(s) && isASCIILower(s[i+1]) {
			i++
			t = append(t, s[i])
		}
	}
	return string(t)
}

// Is this field repeated?
func isRepeated(field *descriptor.FieldDescriptorProto) bool {
	return field.Label != nil && *field.Label == descriptor.FieldDescriptorProto_LABEL_REPEATED
}

// Remove empty part and package part in the name. Join with underscore.
func FullQualifiedTypeNameToGoType(typeName string) string {
	parts := strings.Split(typeName, ".")
	parts = parts[2:]
	for i, v := range parts {
		parts[i] = CamelCase(v)
	}
	return strings.Join(parts, "_")
}

func GoType(field *descriptor.FieldDescriptorProto) string {
	var typ string
	switch *field.Type {
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE:
		typ = "float64"
	case descriptor.FieldDescriptorProto_TYPE_FLOAT:
		typ = "float32"
	case descriptor.FieldDescriptorProto_TYPE_INT64:
		typ = "int64"
	case descriptor.FieldDescriptorProto_TYPE_UINT64:
		typ = "uint64"
	case descriptor.FieldDescriptorProto_TYPE_INT32:
		typ = "int32"
	case descriptor.FieldDescriptorProto_TYPE_UINT32:
		typ = "uint32"
	case descriptor.FieldDescriptorProto_TYPE_FIXED64:
		typ = "uint64"
	case descriptor.FieldDescriptorProto_TYPE_FIXED32:
		typ = "uint32"
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		typ = "bool"
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		typ = "string"
	case descriptor.FieldDescriptorProto_TYPE_GROUP:
		// Group type is deprecated and not supported in proto3. However, Proto3
		// implementations should still be able to parse the group wire format and
		// treat group fields as unknown fields.
		log.Fatal("deprecated FieldDescriptorProto_TYPE_GROUP", field.GetName())
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		typ = FullQualifiedTypeNameToGoType(field.GetTypeName())
	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		typ = "[]byte"
	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		log.Fatal("TODO ", field.GetTypeName())
		// desc := g.ObjectNamed(field.GetTypeName())
		// typ = g.TypeName(desc)
	case descriptor.FieldDescriptorProto_TYPE_SFIXED32:
		typ = "int32"
	case descriptor.FieldDescriptorProto_TYPE_SFIXED64:
		typ = "int64"
	case descriptor.FieldDescriptorProto_TYPE_SINT32:
		typ = "int32"
	case descriptor.FieldDescriptorProto_TYPE_SINT64:
		typ = "int64"
	default:
		log.Fatal("unknown type for", field.GetName())
	}
	if isRepeated(field) {
		typ = "[]" + typ
	}
	return typ
}
