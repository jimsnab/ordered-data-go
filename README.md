# Ordered-Data

This simple library implements an ordered map that has:

* JSON marshal/unmarshal to preserve JSON ordering
* `MustGet*` convenience functions that do some basic type conversions,
  particularly useful when working with JSON data
* Map value extraction
* Replace
* Multi-key removal
* Clone
* Stringify

## Interfaces

`OrderedMap[K comparable, V any]` is a general-purpose ordered map. Create an
ordered map with `NewOrderedMap()` or `NewOrderedMapN()`. Then use its
`Set()` function to assign a value, and `Get()` to retrieve a value. The
rest of the functions in the interface are fairly self-explainitory.

`StringMap` is a convenience alias for `OrderedMap[string, any]`.

## Unmarshaling

To unmarshal an interface in Go, you must first initialize what the interface
points to, and then `json.Unmarshal` will fill it.

Example:
```go
    m := NewOrderedMap[string, any]()
	if err := json.Unmarshal(orderedJson1Json, &m); err != nil {
		panic(err)
	}
```

