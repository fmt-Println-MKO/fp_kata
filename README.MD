[![codecov](https://codecov.io/gh/fmt-Println-MKO/fp_kata/graph/badge.svg?token=CVWSD2NFF0)](https://codecov.io/gh/fmt-Println-MKO/fp_kata)

[![codecov](https://codecov.io/gh/fmt-Println-MKO/fp_kata/graphs/sunburst.svg?token=CVWSD2NFF0)](https://codecov.io/gh/fmt-Println-MKO/fp_kata)

# Functional Programming in Go - Demo Project

![Functional Programming Gopher](docs/dist/theme/img/fp_kata_gopher.png)

This project serves as a **demo application** showcasing how to apply **functional programming patterns** in Go. While Go traditionally follows an imperative, non-functional style, this project explores how concepts like monads, chaining, and folds can be utilized to write cleaner and more expressive code in Go.

The **main branch** contains the classic, non-functional Go implementation. Each step towards functional programming is available in separate branches, showing a gradual transformation of the codebase.

---

## Table of Contents

1. [Project Overview](#project-overview)
2. [Demo Steps](#demo-steps)
3. [Setup and Installation](#setup-and-installation)
4. [Running the Application](#running-the-application)
5. [Generating Code](#generating-code)
6. [Testing and Coverage](#testing-and-coverage)
7. [Project Structure](#project-structure)
8. [Further Reading](#further-reading)
9. [Contributions](#contributions)
10. [License](#license)

---

## Project Overview

This project demonstrates the evolution from **non-functional Go** to **functional Go** programming style. It follows a step-by-step approach, with each step available as a separate branch for clarity. The transformations include:

1. **Classic Non-Functional Style** (Main Branch): The baseline implementation using Go's traditional imperative style.
2. **Step 1 – Applying Monad:** The `OrdersStorage` is refactored to leverage a Monad, encapsulating state and behavior in a composable way.
3. **Step 2 – Chaining Calls:** Functional chaining is introduced, improving readability and reducing boilerplate.
4. **Step 3 – Simplification with Fold:** Code is further simplified by introducing the `fold` operation, which consolidates transformations.
5. **Step 4 – Full Functional Refactor:** The remaining parts of the codebase are refactored to fully adopt functional programming patterns.

For detailed explanations of each step and additional instructions, visit the [project's GitHub page](https://fmt-println-mko.github.io/fp_kata/).

---

## Demo Steps

### Branches Overview:

- **Main Branch**: Classic non-functional implementation.
- **Step 1 Branch**: Introduction of Monads with `OrdersStorage`.
- **Step 2 Branch**: Functional chaining of calls and operations.
- **Step 3 Branch**: Code simplification using the `fold` operation.
- **Step 4 Branch**: Full adoption of functional programming patterns.

---

## Setup and Installation

### Prerequisites:

1. **Go SDK 1.23.5:** Ensure Go is installed. Download it [here](https://go.dev/dl/).
2. **Mockery:** Install Mockery for generating mocks:
   ```shell
   go install github.com/vektra/mockery/v2@latest
   ```
3. **Wire:** Install Wire for dependency injection:
   ```shell
   go install github.com/google/wire/cmd/wire@latest
   ```

---

## Running the Application

To run the demo application, navigate to the desired branch and use the following command:
```shell
go run ./cmd
```

---

## Generating Code

### 1. Generate Dependency Injection

Use `Wire` to auto-generate dependency injection code:
```shell
cd internal/app && go generate ./
```

### 2. Generate Mocks

Generate mock implementations for interfaces with `Mockery`:
```shell
mockery --all --output=mocks
```

---

## Testing and Coverage

Run tests and generate a coverage report to ensure the codebase's quality:

```shell
# Run tests and generate coverage data
go test -cover -coverpkg=$(go list ./... | grep -v '/mocks' | grep -v '/cmd' | tr '\n' ',') ./... -coverprofile=coverage.out

# View the coverage report in a browser
go tool cover -html=coverage.out
```

---

## Project Structure

```plaintext
├── cmd/                     # Entry point for the application
├── internal/                # Application's internal implementation
│   ├── app/                 # Core application logic
│   ├── domain/              # Business entities and logic
│   ├── mocks/               # Auto-generated mocks for testing
│   └── ...                  # Additional packages
└── tests/                   # Unit and integration tests
```

---

## Further Reading

For step-by-step explanations, detailed instructions, and discussions on functional programming in Go, please refer to the [official project documentation](https://fmt-println-mko.github.io/fp_kata/).

---

## Contributions

As this is a demo application, contributions are welcome! If you’d like to contribute:

1. Fork the repository.
2. Create a feature branch:
   ```shell
   git checkout -b feature/<feature-name>
   ```
3. Commit your changes:
   ```shell
   git commit -am 'Add new feature: <feature-description>'
   ```
4. Push to the branch:
   ```shell
   git push origin feature/<feature-name>
   ```
5. Open a Pull Request.

---

## License

This project is licensed under the **Mozilla Public License 2.0** (MPL-2.0).

Under this license:
- You are free to use, modify, and distribute this project.
- You **must retain attribution** to the original author in all distributed copies or modified code.
- Modifications to the source code must also be made available under the same license.

For more details, see the [LICENSE](./LICENSE) file or read the full text of the [MPL-2.0 here](https://www.mozilla.org/en-US/MPL/2.0/).
