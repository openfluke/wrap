package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"syscall/js"
	"time"

	"github.com/openfluke/paragon/v3"
)

// methodWrapper dynamically wraps each method to expose it to JavaScript
func methodWrapper(network interface{}, methodName string) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		method := reflect.ValueOf(network).MethodByName(methodName)
		if !method.IsValid() {
			return fmt.Sprintf("Method %s not found", methodName)
		}

		methodType := method.Type()
		if len(args) == 0 || args[0].IsUndefined() || args[0].String() == "" {
			if methodType.NumIn() == 0 {
				inputs := []reflect.Value{}
				results := method.Call(inputs)
				return serializeResults(results)
			}
			return "No arguments provided"
		}

		var params []interface{}
		paramJSON := args[0].String()
		if err := json.Unmarshal([]byte(paramJSON), &params); err != nil {
			return fmt.Sprintf("Invalid JSON input: %v", err)
		}

		expectedParams := methodType.NumIn()
		if len(params) != expectedParams {
			return fmt.Sprintf("Expected %d parameters, got %d", expectedParams, len(params))
		}

		inputs := make([]reflect.Value, expectedParams)
		for i := 0; i < expectedParams; i++ {
			param := params[i]
			expectedType := methodType.In(i)

			converted, err := convertParameter(param, expectedType, i)
			if err != nil {
				return err.Error()
			}
			inputs[i] = converted
		}

		results := method.Call(inputs)
		return serializeResults(results)
	})
}

// serializeResults converts method results to a JSON string
// convertParameter converts a JavaScript parameter to the expected Go type
func convertParameter(param interface{}, expectedType reflect.Type, paramIndex int) (reflect.Value, error) {
	switch expectedType.Kind() {
	case reflect.Slice:
		return convertSlice(param, expectedType, paramIndex)
	case reflect.Map:
		return convertMap(param, expectedType, paramIndex)
	case reflect.Int:
		if val, ok := param.(float64); ok {
			return reflect.ValueOf(int(val)), nil
		}
		return reflect.Value{}, fmt.Errorf("parameter %d: expected int, got %T", paramIndex, param)
	case reflect.Int8:
		if val, ok := param.(float64); ok {
			return reflect.ValueOf(int8(val)), nil
		}
		return reflect.Value{}, fmt.Errorf("parameter %d: expected int8, got %T", paramIndex, param)
	case reflect.Int16:
		if val, ok := param.(float64); ok {
			return reflect.ValueOf(int16(val)), nil
		}
		return reflect.Value{}, fmt.Errorf("parameter %d: expected int16, got %T", paramIndex, param)
	case reflect.Int32:
		if val, ok := param.(float64); ok {
			return reflect.ValueOf(int32(val)), nil
		}
		return reflect.Value{}, fmt.Errorf("parameter %d: expected int32, got %T", paramIndex, param)
	case reflect.Int64:
		if val, ok := param.(float64); ok {
			return reflect.ValueOf(int64(val)), nil
		}
		return reflect.Value{}, fmt.Errorf("parameter %d: expected int64, got %T", paramIndex, param)
	case reflect.Uint:
		if val, ok := param.(float64); ok {
			return reflect.ValueOf(uint(val)), nil
		}
		return reflect.Value{}, fmt.Errorf("parameter %d: expected uint, got %T", paramIndex, param)
	case reflect.Uint8:
		if val, ok := param.(float64); ok {
			return reflect.ValueOf(uint8(val)), nil
		}
		return reflect.Value{}, fmt.Errorf("parameter %d: expected uint8, got %T", paramIndex, param)
	case reflect.Uint16:
		if val, ok := param.(float64); ok {
			return reflect.ValueOf(uint16(val)), nil
		}
		return reflect.Value{}, fmt.Errorf("parameter %d: expected uint16, got %T", paramIndex, param)
	case reflect.Uint32:
		if val, ok := param.(float64); ok {
			return reflect.ValueOf(uint32(val)), nil
		}
		return reflect.Value{}, fmt.Errorf("parameter %d: expected uint32, got %T", paramIndex, param)
	case reflect.Uint64:
		if val, ok := param.(float64); ok {
			return reflect.ValueOf(uint64(val)), nil
		}
		return reflect.Value{}, fmt.Errorf("parameter %d: expected uint64, got %T", paramIndex, param)
	case reflect.Float32:
		if val, ok := param.(float64); ok {
			return reflect.ValueOf(float32(val)), nil
		}
		return reflect.Value{}, fmt.Errorf("parameter %d: expected float32, got %T", paramIndex, param)
	case reflect.Float64:
		if val, ok := param.(float64); ok {
			return reflect.ValueOf(val), nil
		}
		return reflect.Value{}, fmt.Errorf("parameter %d: expected float64, got %T", paramIndex, param)
	case reflect.Bool:
		if val, ok := param.(bool); ok {
			return reflect.ValueOf(val), nil
		}
		return reflect.Value{}, fmt.Errorf("parameter %d: expected bool, got %T", paramIndex, param)
	case reflect.String:
		if val, ok := param.(string); ok {
			return reflect.ValueOf(val), nil
		}
		return reflect.Value{}, fmt.Errorf("parameter %d: expected string, got %T", paramIndex, param)
	default:
		if expectedType == reflect.TypeOf(time.Duration(0)) {
			if val, ok := param.(float64); ok {
				return reflect.ValueOf(time.Duration(val)), nil
			}
			return reflect.Value{}, fmt.Errorf("parameter %d: expected duration, got %T", paramIndex, param)
		}
		return reflect.Zero(expectedType), fmt.Errorf("parameter %d: unsupported type %s", paramIndex, expectedType.String())
	}
}

// convertSlice handles slice type conversions including multidimensional slices
func convertSlice(param interface{}, expectedType reflect.Type, paramIndex int) (reflect.Value, error) {
	val, ok := param.([]interface{})
	if !ok {
		return reflect.Value{}, fmt.Errorf("parameter %d: expected slice, got %T", paramIndex, param)
	}

	elemType := expectedType.Elem()
	slice := reflect.MakeSlice(expectedType, len(val), len(val))

	for j, v := range val {
		switch elemType.Kind() {
		case reflect.Int:
			if num, ok := v.(float64); ok {
				slice.Index(j).SetInt(int64(num))
			} else {
				return reflect.Value{}, fmt.Errorf("parameter %d: invalid slice element type %T at index %d", paramIndex, v, j)
			}
		case reflect.Int8:
			if num, ok := v.(float64); ok {
				slice.Index(j).SetInt(int64(num))
			} else {
				return reflect.Value{}, fmt.Errorf("parameter %d: invalid slice element type %T at index %d", paramIndex, v, j)
			}
		case reflect.Int16:
			if num, ok := v.(float64); ok {
				slice.Index(j).SetInt(int64(num))
			} else {
				return reflect.Value{}, fmt.Errorf("parameter %d: invalid slice element type %T at index %d", paramIndex, v, j)
			}
		case reflect.Int32:
			if num, ok := v.(float64); ok {
				slice.Index(j).SetInt(int64(num))
			} else {
				return reflect.Value{}, fmt.Errorf("parameter %d: invalid slice element type %T at index %d", paramIndex, v, j)
			}
		case reflect.Int64:
			if num, ok := v.(float64); ok {
				slice.Index(j).SetInt(int64(num))
			} else {
				return reflect.Value{}, fmt.Errorf("parameter %d: invalid slice element type %T at index %d", paramIndex, v, j)
			}
		case reflect.Uint:
			if num, ok := v.(float64); ok {
				slice.Index(j).SetUint(uint64(num))
			} else {
				return reflect.Value{}, fmt.Errorf("parameter %d: invalid slice element type %T at index %d", paramIndex, v, j)
			}
		case reflect.Uint8:
			if num, ok := v.(float64); ok {
				slice.Index(j).SetUint(uint64(num))
			} else {
				return reflect.Value{}, fmt.Errorf("parameter %d: invalid slice element type %T at index %d", paramIndex, v, j)
			}
		case reflect.Uint16:
			if num, ok := v.(float64); ok {
				slice.Index(j).SetUint(uint64(num))
			} else {
				return reflect.Value{}, fmt.Errorf("parameter %d: invalid slice element type %T at index %d", paramIndex, v, j)
			}
		case reflect.Uint32:
			if num, ok := v.(float64); ok {
				slice.Index(j).SetUint(uint64(num))
			} else {
				return reflect.Value{}, fmt.Errorf("parameter %d: invalid slice element type %T at index %d", paramIndex, v, j)
			}
		case reflect.Uint64:
			if num, ok := v.(float64); ok {
				slice.Index(j).SetUint(uint64(num))
			} else {
				return reflect.Value{}, fmt.Errorf("parameter %d: invalid slice element type %T at index %d", paramIndex, v, j)
			}
		case reflect.Float32:
			if num, ok := v.(float64); ok {
				slice.Index(j).SetFloat(num)
			} else {
				return reflect.Value{}, fmt.Errorf("parameter %d: invalid slice element type %T at index %d", paramIndex, v, j)
			}
		case reflect.Float64:
			if num, ok := v.(float64); ok {
				slice.Index(j).SetFloat(num)
			} else {
				return reflect.Value{}, fmt.Errorf("parameter %d: invalid slice element type %T at index %d", paramIndex, v, j)
			}
		case reflect.Bool:
			if b, ok := v.(bool); ok {
				slice.Index(j).SetBool(b)
			} else {
				return reflect.Value{}, fmt.Errorf("parameter %d: invalid slice element type %T at index %d", paramIndex, v, j)
			}
		case reflect.String:
			if s, ok := v.(string); ok {
				slice.Index(j).SetString(s)
			} else {
				return reflect.Value{}, fmt.Errorf("parameter %d: invalid slice element type %T at index %d", paramIndex, v, j)
			}
		case reflect.Slice:
			// Handle 2D/3D slices recursively
			converted, err := convertSlice(v, elemType, paramIndex)
			if err != nil {
				return reflect.Value{}, err
			}
			slice.Index(j).Set(converted)
		case reflect.Struct:
			// Handle struct slices (like layer size structs)
			if structVal, ok := v.(map[string]interface{}); ok {
				structValue := reflect.New(elemType).Elem()
				for fieldName, fieldVal := range structVal {
					field := structValue.FieldByName(fieldName)
					if field.IsValid() && field.CanSet() {
						if num, ok := fieldVal.(float64); ok && field.Kind() == reflect.Int {
							field.SetInt(int64(num))
						}
					}
				}
				slice.Index(j).Set(structValue)
			}
		default:
			return reflect.Value{}, fmt.Errorf("parameter %d: unsupported slice element type %s", paramIndex, elemType.String())
		}
	}

	return slice, nil
}

// convertMap handles map type conversions
func convertMap(param interface{}, expectedType reflect.Type, paramIndex int) (reflect.Value, error) {
	jsonMap, ok := param.(map[string]interface{})
	if !ok {
		return reflect.Value{}, fmt.Errorf("parameter %d: expected map, got %T", paramIndex, param)
	}

	keyType := expectedType.Key()
	valueType := expectedType.Elem()
	mapValue := reflect.MakeMap(expectedType)

	for keyStr, val := range jsonMap {
		var keyValue reflect.Value
		var err error

		switch keyType.Kind() {
		case reflect.Int:
			key, parseErr := strconv.Atoi(keyStr)
			if parseErr != nil {
				return reflect.Value{}, fmt.Errorf("parameter %d: invalid map key %s", paramIndex, keyStr)
			}
			keyValue = reflect.ValueOf(key)
		case reflect.String:
			keyValue = reflect.ValueOf(keyStr)
		default:
			return reflect.Value{}, fmt.Errorf("parameter %d: unsupported map key type %s", paramIndex, keyType.String())
		}

		valueValue, err := convertParameter(val, valueType, paramIndex)
		if err != nil {
			return reflect.Value{}, err
		}

		mapValue.SetMapIndex(keyValue, valueValue)
	}

	return mapValue, nil
}

// serializeResults converts method results to JSON
func serializeResults(results []reflect.Value) string {
	if len(results) == 0 {
		return "[]"
	}

	output := make([]interface{}, len(results))
	for i, result := range results {
		output[i] = result.Interface()
	}

	resultJSON, err := json.Marshal(output)
	if err != nil {
		return fmt.Sprintf("Failed to marshal results: %v", err)
	}
	return string(resultJSON)
}

// newParagonWrapper creates a wrapper based on the numeric type string
func newParagonWrapper(numericType string) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// Parse the NewNetwork arguments
		if len(args) < 3 {
			errMsg := "Expected 3 arguments: layerSizes, activations, fullyConnected"
			js.Global().Call("console", "error", errMsg)
			return errMsg
		}

		var layerSizes []struct{ Width, Height int }
		var activations []string
		var fullyConnected []bool

		if err := json.Unmarshal([]byte(args[0].String()), &layerSizes); err != nil {
			errMsg := fmt.Sprintf("Invalid layerSizes JSON: %v", err)
			js.Global().Call("console", "error", errMsg)
			return errMsg
		}
		if err := json.Unmarshal([]byte(args[1].String()), &activations); err != nil {
			errMsg := fmt.Sprintf("Invalid activations JSON: %v", err)
			js.Global().Call("console", "error", errMsg)
			return errMsg
		}
		if err := json.Unmarshal([]byte(args[2].String()), &fullyConnected); err != nil {
			errMsg := fmt.Sprintf("Invalid fullyConnected JSON: %v", err)
			js.Global().Call("console", "error", errMsg)
			return errMsg
		}

		var network interface{}
		var err error

		// Create network based on the numeric type
		switch numericType {
		case "float32":
			network, err = paragon.NewNetwork[float32](layerSizes, activations, fullyConnected)
		case "float64":
			network, err = paragon.NewNetwork[float64](layerSizes, activations, fullyConnected)
		case "int":
			network, err = paragon.NewNetwork[int](layerSizes, activations, fullyConnected)
		case "int8":
			network, err = paragon.NewNetwork[int8](layerSizes, activations, fullyConnected)
		case "int16":
			network, err = paragon.NewNetwork[int16](layerSizes, activations, fullyConnected)
		case "int32":
			network, err = paragon.NewNetwork[int32](layerSizes, activations, fullyConnected)
		case "int64":
			network, err = paragon.NewNetwork[int64](layerSizes, activations, fullyConnected)
		case "uint":
			network, err = paragon.NewNetwork[uint](layerSizes, activations, fullyConnected)
		case "uint8":
			network, err = paragon.NewNetwork[uint8](layerSizes, activations, fullyConnected)
		case "uint16":
			network, err = paragon.NewNetwork[uint16](layerSizes, activations, fullyConnected)
		case "uint32":
			network, err = paragon.NewNetwork[uint32](layerSizes, activations, fullyConnected)
		case "uint64":
			network, err = paragon.NewNetwork[uint64](layerSizes, activations, fullyConnected)
		default:
			errMsg := fmt.Sprintf("Unsupported numeric type: %s", numericType)
			js.Global().Call("console", "error", errMsg)
			return errMsg
		}

		if err != nil {
			errMsg := fmt.Sprintf("Error: failed to create network: %v", err)
			js.Global().Call("console", "error", errMsg)
			return errMsg
		}

		obj := js.Global().Get("Object").New()

		// Use reflection to get all methods
		networkValue := reflect.ValueOf(network)
		networkType := networkValue.Type()

		for i := 0; i < networkType.NumMethod(); i++ {
			method := networkType.Method(i)
			// Only export public methods
			if method.Name[0] >= 'A' && method.Name[0] <= 'Z' {
				obj.Set(method.Name, methodWrapper(network, method.Name))
			}
		}

		return obj
	})
}

func main() {
	// Register factory functions for each type
	js.Global().Set("NewNetworkFloat32", newParagonWrapper("float32"))
	js.Global().Set("NewNetworkFloat64", newParagonWrapper("float64"))
	js.Global().Set("NewNetworkInt", newParagonWrapper("int"))
	js.Global().Set("NewNetworkInt8", newParagonWrapper("int8"))
	js.Global().Set("NewNetworkInt16", newParagonWrapper("int16"))
	js.Global().Set("NewNetworkInt32", newParagonWrapper("int32"))
	js.Global().Set("NewNetworkInt64", newParagonWrapper("int64"))
	js.Global().Set("NewNetworkUint", newParagonWrapper("uint"))
	js.Global().Set("NewNetworkUint8", newParagonWrapper("uint8"))
	js.Global().Set("NewNetworkUint16", newParagonWrapper("uint16"))
	js.Global().Set("NewNetworkUint32", newParagonWrapper("uint32"))
	js.Global().Set("NewNetworkUint64", newParagonWrapper("uint64"))

	// Default to float32
	js.Global().Set("NewNetwork", newParagonWrapper("float32"))

	select {}
}
