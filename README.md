# Paragon WASM Wrapper

## Overview

This project provides a WebAssembly (WASM) wrapper for the Paragon AI framework, written in Go. It enables the use of Paragon's neural network capabilities in JavaScript environments, such as web browsers, by compiling Go code to WASM. The wrapper exposes Paragon's network creation and methods to JavaScript, supporting multiple numeric types for flexibility.

## Git clone from

```
go get github.com/openfluke/paragon/v3@v3.0.0
```

```
go get github.com/openfluke/webgpu@ad2e76f
```

## Purpose

The main goal is to bridge the Paragon AI framework with JavaScript applications, allowing developers to leverage Paragon's neural network functionality in web-based projects. This is achieved by:

- Compiling the Go code to WASM using `GOOS=js GOARCH=wasm go build -o paragon.wasm main.go`.
- Providing a dynamic interface to create and interact with Paragon neural networks from JavaScript.

## Features

- **Dynamic Method Exposure**: Uses Go's reflection to expose public methods of the Paragon network to JavaScript, allowing seamless method calls.
- **Type Flexibility**: Supports multiple numeric types (`float32`, `float64`, `int`, `int8`, `int16`, `int32`, `int64`, `uint`, `uint8`, `uint16`, `uint32`, `uint64`) for network creation.
- **JSON Parameter Handling**: Accepts JSON-encoded parameters from JavaScript, converting them to appropriate Go types for Paragon's methods.
- **Slice and Map Support**: Handles complex data structures like slices (including multidimensional) and maps, converting them between JavaScript and Go.
- **Error Handling**: Provides detailed error messages for invalid inputs, type mismatches, or JSON parsing issues.

## Performance

The following table summarizes the performance of various activation functions tested with the Paragon WASM wrapper using `test.html`. The tests measure execution time (in milliseconds) on CPU and GPU, calculate the speedup (CPU time / GPU time), and verify result consistency (`Match`).

| Type    | Activation | CPU (ms) | GPU (ms) | Speedup | Match |
| ------- | ---------- | -------- | -------- | ------- | ----- |
| float32 | linear     | 1102.90  | 1090.70  | 1.01    | ✅    |
| float32 | relu       | 1137.50  | 1095.00  | 1.04    | ✅    |
| float32 | leaky_relu | 1139.70  | 1091.60  | 1.04    | ✅    |
| float32 | elu        | 1105.60  | 1082.90  | 1.02    | ✅    |
| float32 | swish      | 1106.30  | 1097.40  | 1.01    | ✅    |
| float32 | gelu       | 1050.80  | 1116.30  | 0.94    | ✅    |
| float32 | tanh       | 1033.90  | 1102.60  | 0.94    | ✅    |
| float32 | softmax    | 1045.70  | 1158.30  | 0.90    | ✅    |
| int32   | linear     | 1043.10  | 1101.20  | 0.95    | ✅    |
| int32   | relu       | 1036.40  | 1034.30  | 1.00    | ✅    |
| int32   | leaky_relu | 1147.60  | 1148.70  | 1.00    | ✅    |
| int32   | elu        | 1061.10  | 1077.10  | 0.99    | ✅    |
| int32   | swish      | 1163.70  | 1065.00  | 1.09    | ✅    |
| int32   | gelu       | 1083.10  | 1167.90  | 0.93    | ✅    |
| int32   | tanh       | 1162.00  | 1065.40  | 1.09    | ✅    |
| int32   | softmax    | 1068.10  | 1098.20  | 0.97    | ✅    |
| uint32  | linear     | 1228.90  | 1048.00  | 1.17    | ✅    |
| uint32  | relu       | 1054.20  | 1060.80  | 0.99    | ✅    |
| uint32  | leaky_relu | 1092.80  | 1141.30  | 0.96    | ✅    |
| uint32  | elu        | 1052.50  | 1039.70  | 1.01    | ✅    |
| uint32  | swish      | 1050.40  | 1030.90  | 1.02    | ✅    |
| uint32  | gelu       | 1115.60  | 1171.50  | 0.95    | ✅    |
| uint32  | tanh       | 1180.60  | 1068.40  | 1.11    | ✅    |
| uint32  | softmax    | 1088.00  | 1055.20  | 1.03    | ✅    |

## How It Works

The Go code in `main.go` does the following:

1. **Network Creation**: The `newParagonWrapper` function creates a Paragon neural network based on a specified numeric type, using parameters like layer sizes, activations, and connectivity options passed as JSON from JavaScript.
2. **Method Wrapping**: The `methodWrapper` function dynamically wraps Paragon's public methods, enabling JavaScript to call them. It handles parameter conversion from JSON to Go types and serializes results back to JSON.
3. **Type Conversion**: Functions like `convertParameter`, `convertSlice`, and `convertMap` convert JavaScript inputs to Go types, supporting primitives, slices, maps, and structs.
4. **WASM Integration**: The `main` function registers factory functions (`NewNetworkFloat32`, `NewNetworkFloat64`, etc.) in the JavaScript global scope, with `NewNetwork` defaulting to `float32`. The program runs indefinitely using `select {}` to keep the WASM module active.

## Compilation

To compile the Go code to WebAssembly:

```bash
GOOS=js GOARCH=wasm go build -o paragon.wasm main.go
```

This generates `paragon.wasm`, which can be loaded in a JavaScript environment using the WebAssembly API. You'll also need the `wasm_exec.js` file provided by Go to instantiate the WASM module.

## Hosting

To serve the WASM module and associated files (e.g., `test.html`) locally for testing:

```bash
python -m http.server 8000
```

This starts a simple HTTP server on port 8000. Access the application by navigating to `http://localhost:8000` in a web browser.

## Usage

1. **Include WASM in JavaScript**:

   - Load `wasm_exec.js` and `paragon.wasm` in your web application.
   - Instantiate the WASM module using Go's runtime.

2. **Create a Network**:

   ```javascript
   const layerSizes = [
     { Width: 2, Height: 1 },
     { Width: 3, Height: 1 },
   ];
   const activations = ["relu", "sigmoid"];
   const fullyConnected = [true];
   const network = await NewNetwork(
     JSON.stringify(layerSizes),
     JSON.stringify(activations),
     JSON.stringify(fullyConnected)
   );
   ```

3. **Call Methods**:

   ```javascript
   const result = await network.SomeMethod(JSON.stringify([1.0, 2.0]));
   console.log(result); // JSON string with method results
   ```

4. **Supported Types**:
   - Use `NewNetworkFloat32`, `NewNetworkInt`, etc., to create networks with specific numeric types.
   - Pass parameters as JSON arrays or objects, which are converted to Go types like slices, maps, or structs.

## Dependencies

- **Go**: Version compatible with WASM compilation (e.g., Go 1.21+).
- **Paragon**: The Paragon AI framework (`github.com/openfluke/paragon`) for neural network functionality.
- **JavaScript Environment**: A browser or Node.js with WebAssembly support.
- **Python**: Required for running the local HTTP server (`python -m http.server 8000`).

## Limitations

- **Numeric Types**: Only predefined numeric types are supported. Custom types require additional wrapper functions.
- **Complex Structs**: Limited support for nested structs; only simple struct fields (e.g., `Width`, `Height`) are handled.
- **Performance**: Reflection and JSON parsing may introduce overhead for high-performance applications.
- **Browser Compatibility**: Ensure the target environment supports WebAssembly.

## Contributing

Contributions are welcome! Please submit issues or pull requests to improve the wrapper, add support for more types, or optimize performance.
