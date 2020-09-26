package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"regexp"
	"strings"
	"text/template"
)

type definition struct {
	packageName string
	objects     map[string]nodeStr
}

type nodeStr struct {
	nodeType    string              //Истинный тип структуры (присутствует только в базовых типах и их обёртках)
	fieldList   map[string]fieldStr //Список полей структуры
	isValidated bool                //Нужно ли писать валидаторы
}

type fieldStr struct {
	primaryFieldType   string //Основной тип данных
	secondaryFieldType string //Второстепенный тип данных для мапов
	fieldTag           string //Тег для валидации
	isList             bool   //Является ли списком
}

type conditionParams struct {
	Obj       string
	Value     string
	FieldName string
	Condition string
}

func main() {
	filename := os.Args[1]

	definitionStr, err := analyzeDeclaration(filename)
	if err != nil {
		fmt.Errorf("%w", err)
	}
	err = generateValidation(filename, definitionStr)
	if err != nil {
		fmt.Errorf("%w", err)
	}
	return
}

func analyzeDeclaration(filename string) (*definition, error) {
	var re = regexp.MustCompile(`^.*validate:"(.*)".*$`)
	definitionStr := definition{}
	definitionStr.objects = map[string]nodeStr{}
	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	ast.Inspect(file, func(x ast.Node) bool {
		switch nodeType := x.(type) {
		case *ast.File:
			definitionStr.packageName = nodeType.Name.Name
		case *ast.GenDecl:
			for _, spec := range nodeType.Specs {
				switch specType := spec.(type) {
				case *ast.TypeSpec:
					currentNode := nodeStr{}
					currentNode.fieldList = map[string]fieldStr{}
					FieldType, ok := specType.Type.(*ast.Ident)
					if ok {
						currentNode.nodeType = FieldType.Name
					}
					StructType, ok := specType.Type.(*ast.StructType)
					if ok {
						for _, field := range StructType.Fields.List {
							currentNode.fieldList[field.Names[0].Name] = analyzeField(field, re)
						}
					}
					definitionStr.objects[specType.Name.String()] = currentNode
				default:
				}
			}
		default:
			return true
		}
		return true
	})

	analyzeStruct(definitionStr)
	return &definitionStr, nil
}

func analyzeField(declaredField *ast.Field, re *regexp.Regexp) (currentField fieldStr) {
	currentField.isList = false
	mapType, ok := declaredField.Type.(*ast.MapType)
	if ok {
		currentField.primaryFieldType = mapType.Value.(*ast.Ident).String()
		currentField.secondaryFieldType = mapType.Key.(*ast.Ident).String()
	}
	arrayType, ok := declaredField.Type.(*ast.ArrayType)
	if ok {
		currentField.primaryFieldType = arrayType.Elt.(*ast.Ident).String()
		currentField.isList = true
	}
	fieldType, ok := declaredField.Type.(*ast.Ident)
	if ok {
		currentField.primaryFieldType = fieldType.Name
	}
	if declaredField.Tag != nil {
		isExist := re.MatchString(declaredField.Tag.Value)
		if isExist {
			currentField.fieldTag = re.ReplaceAllString(declaredField.Tag.Value, "$1")
		} else {
			currentField.fieldTag = ""
		}
	}
	return
}

func analyzeStruct(definitionStr definition) {
	for classKey, classValue := range definitionStr.objects {
		for fieldKey, fieldValue := range classValue.fieldList {
			if fieldValue.fieldTag != "" && classValue.isValidated == false {
				classValue.isValidated = true
				definitionStr.objects[classKey] = classValue
			}
			if val, ok := definitionStr.objects[fieldValue.primaryFieldType]; ok {
				if val.nodeType != "" {
					fieldValue.primaryFieldType = val.nodeType
					definitionStr.objects[classKey].fieldList[fieldKey] = fieldValue
				}
			}
		}
	}
}

func generateValidation(filename string, definitionStr *definition) error {
	// Открываем
	var re = regexp.MustCompile(`^(.*).go$`)
	newFileName := re.ReplaceAllString(filename, "$1") + "_validation_generated.go"
	validationFile, err := os.Create(newFileName)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer validationFile.Close()
	funcTmp := template.New("func")

	// Вступление
	t := template.Must(funcTmp.Parse(codeBlocks["main"]["intro"]))
	err = t.Execute(validationFile, conditionParams{definitionStr.packageName, "", "", ""})
	if err != nil {
		return fmt.Errorf("failed to generate validation: %w", err)
	}

	// Сюжет
	for key, value := range definitionStr.objects {
		if value.isValidated {
			// Объявления функций валидации
			t := template.Must(funcTmp.Parse(codeBlocks["main"]["validateDeclaration"]))
			err = t.Execute(validationFile, conditionParams{key, "", "", ""})
			if err != nil {
				return fmt.Errorf("failed to generate validation: %w", err)
			}
			for objKey, objValue := range value.fieldList {
				if objValue.fieldTag != "" {
					conditions := strings.Split(objValue.fieldTag, "|")
					if objValue.isList {
						_, err = validationFile.WriteString("	for _, v := range o." + objKey + " {\n")
						if err != nil {
							return fmt.Errorf("failed to generate validation: %w", err)
						}
						for _, v := range conditions {
							p := conditionParams{
								"v", v[strings.Index(v, ":")+1:], objKey, v,
							}
							if v[:strings.Index(v, ":")] == "in" && objValue.primaryFieldType == "string" {
								p.Value = strings.ReplaceAll(p.Value, ",", "\",\"")
							}
							t := template.Must(funcTmp.Parse(codeBlocks[objValue.primaryFieldType][v[:strings.Index(v, ":")]]))
							err = t.Execute(validationFile, p)
							if err != nil {
								return fmt.Errorf("failed to generate validation: %w", err)
							}
						}
						_, err = validationFile.WriteString("\n\t\tif err.Err != nil {\n\t\t\tbreak\n\t\t}\n\t}\n")
						if err != nil {
							return fmt.Errorf("failed to generate validation: %w", err)
						}
					} else {
						for _, v := range conditions {
							p := conditionParams{
								"o." + objKey, v[strings.Index(v, ":")+1:], objKey, v,
							}
							if v[:strings.Index(v, ":")] == "in" && objValue.primaryFieldType == "string" {
								p.Value = strings.ReplaceAll(p.Value, ",", "\",\"")
							}
							t := template.Must(funcTmp.Parse(codeBlocks[objValue.primaryFieldType][v[:strings.Index(v, ":")]]))
							err = t.Execute(validationFile, p)
							if err != nil {
								return fmt.Errorf("failed to generate validation: %w", err)
							}
						}
					}
				}
			}
			_, err = validationFile.WriteString("\treturn validationErrorList, nil\n}\n")
			if err != nil {
				return fmt.Errorf("failed to generate validation: %w", err)
			}
		}
	}

	return nil
}

var codeBlocks = map[string]map[string]string{
	"string": {
		"len": `
	func() {
		if len({{.Obj}}) != {{.Value}} {
			err = ValidationError{"{{.FieldName}}", errors.New("failed to validate field {{.FieldName}} with condition {{.Condition}}")}
			validationErrorList = append(validationErrorList, err)
		}
	}()
`,
		"regexp": `
	func() {
		re := regexp.MustCompile("{{.Value}}")
		if !re.MatchString({{.Obj}}) {
			err = ValidationError{"{{.FieldName}}", errors.New("failed to validate field {{.FieldName}} with condition {{.Condition}}")}
			validationErrorList = append(validationErrorList, err)
		}
	}()
`,
		"in": `
	func() {
		switch {{.Obj}} {
		case "{{.Value}}":
		default:
			err = ValidationError{"{{.FieldName}}", errors.New("failed to validate field {{.FieldName}} with condition {{.Condition}}")}
			validationErrorList = append(validationErrorList, err)
		}
	}()
`,
	},
	"int": {
		"min": `
	func() {
		if {{.Obj}} < {{.Value}} {
			err = ValidationError{"{{.FieldName}}", errors.New("failed to validate field {{.FieldName}} with condition {{.Condition}}")}
			validationErrorList = append(validationErrorList, err)
		}
	}()
`,
		"max": `
	func() {
		if {{.Obj}} > {{.Value}} {
			err = ValidationError{"{{.FieldName}}", errors.New("failed to validate field {{.FieldName}} with condition {{.Condition}}")}
			validationErrorList = append(validationErrorList, err)
		}
	}()
`,
		"in": `
	func() {
		switch {{.Obj}} {
		case {{.Value}}:
		default:
			err = ValidationError{"{{.FieldName}}", errors.New("failed to validate field {{.FieldName}} with condition {{.Condition}}")}
			validationErrorList = append(validationErrorList, err)
		}
	}()
`,
	},
	"main": {
		"intro": `
// Code generated by cool go-validate tool; DO NOT EDIT.

package {{.Obj}}

import (
	"errors"
	"regexp"
)

type ValidationError struct {
	Field   string
	Err error
}
`,
		"validateDeclaration": `
func (o {{.Obj}}) Validate() ([]ValidationError, error) {
	validationErrorList := []ValidationError{}
	err := ValidationError{}
`,
	},
}
