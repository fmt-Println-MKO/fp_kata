# Functional Programming in Go
## 1. Understanding Functional Programming
### 1.1 Core Principles
- Pure functions
  - Same input always produces same output
  - No side effects
  - Makes code predictable and testable

- Immutability
  - Data doesn't change after creation
  - Creates predictable state management
  - Helps with concurrent programming

- First-class functions
  - Functions can be assigned to variables
  - Functions can be passed as arguments
  - Functions can be returned from other functions

- Higher-order functions
  - Functions that take functions as parameters
  - Functions that return functions
  - Enables composition and reusability

### 1.2 Advanced FP Concepts
- Pattern Matching
  - Destructuring data based on its shape
  - Handling different cases explicitly
  - Makes code more declarative
  - Essential for working with algebraic data types

- Algebraic Data Types (ADTs)
  - Sum types (OR relationships)
  - Product types (AND relationships)
  - Type-safe way to model domain concepts
  - Fundamental for encoding business rules in types

- Generic Methods
  - Methods that work with any type
  - Enable type-safe abstractions
  - Crucial for building reusable FP abstractions like functors and monads

- Type Classes
  - Define behavior for types
  - Enable ad-hoc polymorphism
  - Allow for powerful abstractions

## 2. FP Features Not Natively Available in Go
### 2.1 Pattern Matching
#### Haskell Example
```haskell
data Shape = Circle Float | Rectangle Float Float

area :: Shape -> Float
area (Circle r) = pi * r * r
area (Rectangle w h) = w * h
```

#### Scala Example
```scala
sealed trait Shape
case class Circle(radius: Float) extends Shape
case class Rectangle(width: Float, height: Float) extends Shape

def area(shape: Shape): Float = shape match {
  case Circle(r) => math.Pi * r * r
  case Rectangle(w, h) => w * h
}
```

#### Go Workarounds
```go
type Shape interface {
    Area() float64
}

type Circle struct {
    radius float64
}

type Rectangle struct {
    width, height float64
}

// Type-specific implementation instead of pattern matching
func (c Circle) Area() float64 {
    return math.Pi * c.radius * c.radius
}

func (r Rectangle) Area() float64 {
    return r.width * r.height
}

// Alternative using type switches
func CalculateArea(s Shape) float64 {
    switch v := s.(type) {
    case Circle:
        return math.Pi * v.radius * v.radius
    case Rectangle:
        return v.width * v.height
    default:
        return 0
    }
}
```

### 2.2 Algebraic Data Types
#### OCaml Example
```ocaml
type result =
  | Success of string
  | Error of string
  | NotFound

let handle_result r = match r with
  | Success msg -> print_endline ("Success: " ^ msg)
  | Error err -> print_endline ("Error: " ^ err)
  | NotFound -> print_endline "Not found"
```

#### F# Example
```fsharp
type ValidationResult =
    | Valid of string
    | Invalid of string list

let validate result =
    match result with
    | Valid s -> printfn "Valid: %s" s
    | Invalid errors ->
        errors
        |> List.iter (printfn "Error: %s")
```

#### Go Workarounds
```go
// Using interfaces and structs to simulate sum types
type ValidationResult interface {
    isValidationResult()
}

type Valid struct {
    value string
}

type Invalid struct {
    errors []string
}

func (Valid) isValidationResult()   {}
func (Invalid) isValidationResult() {}

// Using type assertions for handling
func handleValidation(r ValidationResult) {
    switch v := r.(type) {
    case Valid:
        fmt.Printf("Valid: %s\n", v.value)
    case Invalid:
        for _, err := range v.errors {
            fmt.Printf("Error: %s\n", err)
        }
    }
}
```

### 2.3 Method Generics
#### Scala Example
```scala
class Container[A](value: A) {
    def map[B](f: A => B): Container[B] =
        new Container(f(value))
}
```

#### Haskell Example
```haskell
class Functor f where
    fmap :: (a -> b) -> f a -> f b

instance Functor Maybe where
    fmap _ Nothing = Nothing
    fmap f (Just x) = Just (f x)
```

#### Go Limitations and Workarounds
```go
// Instead of method generics, we need separate methods or type-specific implementations
type StringContainer struct {
    value string
}

type IntContainer struct {
    value int
}

// We can't write a generic map method that returns a different type
// Instead we need specific methods:
func (c StringContainer) MapToInt(f func(string) int) IntContainer {
    return IntContainer{f(c.value)}
}

// Or use generic functions instead of methods
func Map[T, U any](container Container[T], f func(T) U) Container[U] {
    return Container[U]{f(container.value)}
}

// This means we lose method chaining capability:
// Instead of: container.Map(f1).Map(f2)
// We must write: Map(Map(container, f1), f2)
```

## 3. FP Features Available in Go
### 3.1 First-Class Functions
- Example in Haskell
- Example in Scala
- Implementation in Go
- Common use cases

### 3.2 Closures
- Example in Functional Languages
- Go implementation
- Practical applications

### 3.3 Higher-Order Functions
- Examples from functional languages
- Go implementation patterns
- Best practices

## 4. Functional Programming Libraries in Go
### 4.1 samber/mo: Monads and FP patterns
- Overview of available monads
- Result monad
- Option monad
- Common use patterns
- Integration with existing Go code

### 4.2 samber/lo: Functional helpers
- Available operations
- Map, Filter, Reduce implementations
- Performance considerations
- Integration examples

## 5. Practical Examples
### 5.1 Error Handling with Result Monad
- Traditional Go error handling
```go
data, err := someFunction()
if err != nil {
    return nil, err
}
```
- Result monad approach
```go
result := someFunction()
    .Map(func(data Type) Type {
        // transform data
    })
    .MapError(func(err error) error {
        // transform error
    })
```
- Benefits and tradeoffs

### 5.2 Working with Optional Values
- Traditional nil checking
- Option monad usage
- Preventing nil pointer exceptions
- Cleaner code examples

### 5.3 Collection Operations
- Traditional Go loops vs functional approach
- Filter examples
- Map usage
- FlatMap scenarios
- Working around method generic limitations

## 6. Real-World Applications
### 6.1 Data Processing Pipelines
- Building composable operations
- Error handling in pipelines
- Testing strategies

### 6.2 API Response Handling
- Using Result for HTTP responses
- Composing multiple API calls
- Error propagation

## 7. Benefits and Trade-offs
### 7.1 Advantages
- Improved code organization
- Better testability
- Reduced side effects
- Clearer intent
- Reusable components

### 7.2 Challenges
- Learning curve
- Go language limitations
- Performance considerations
- Team adoption

## 8. Conclusion
- FP is possible in Go
- Benefits of functional thinking
- Practical adoption strategy
- Resources for further learning
