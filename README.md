# go-html-retriever

A lightweight Go utility for parsing HTML and querying nodes by id, tag, and class.

It is built on top of `golang.org/x/net/html` and provides small wrapper types that make traversal and extraction easier.

## Quick Start

```go
package main

import (
    "fmt"
    "strings"

    "github.com/we1sper/go-html-retriever"
)

func main() {
    html := `
    <html>
        <body>
            <div id="main" class="container layout">
                <h1 class="title">Hello World!</h1>
                <p class="paragraph">A Simple Page</p>
                <p class="paragraph">Used for testing</p>
                <p class="paragraph tail">Created by welsper</p>
            </div>
        </body>
    </html>`
	
    root, err := htmlretriever.From(strings.NewReader(html))
    if err != nil {
        panic(err)
    }

    m := root.GetElementById("main")
    if !m.Present() {
        fmt.Println("main container not found")
        return
    }

    title := m.GetElementsByClassName("title").First()
    if text, ok := title.GetText(); ok {
        fmt.Printf("title: %s\n", text)
    }

    paragraphs := m.GetElementsByTagName("p")
    fmt.Printf("paragraph count: %d\n", paragraphs.Len())
	
	author := m.GetElementsByClassName("paragraph", "tail").First().Unsafe().GetText()
    fmt.Printf("author: %s\n", author)
}

```

## Core Types

- `EnhancedNode`: wrapper around a single `*html.Node`.
- `EnhancedNodes`: wrapper around a slice of `*EnhancedNode`.
- `UnsafeNode`: helper that skips `(value, present)` checks and returns zero values when missing.

## API Overview

### Parse

- `From(reader io.Reader) (*EnhancedNode, error)` reads HTML from the provided reader and returns the root node wrapped in an `EnhancedNode`.

### Query

On `*EnhancedNode`:

- `GetElementById(id string) *EnhancedNode` returns the first node with the specified id.
- `GetElementsByTagName(tag string) *EnhancedNodes` returns all nodes with the specified tag name.
- `GetElementsByClassName(classes ...string) *EnhancedNodes` returns all nodes that contain all specified classes.
- `Filter(filter func(*EnhancedNode) bool) *EnhancedNodes` returns all nodes in the subtree that satisfy the filter function.

On `*EnhancedNodes`:

- `GetElementById(id string) *EnhancedNode` returns the first node with the specified id.
- `GetElementsByTagName(tag string) *EnhancedNodes` returns all nodes with the specified tag name.
- `GetElementsByClassName(classes ...string) *EnhancedNodes` returns all nodes that contain all specified classes.
- `Filter(filter func(*EnhancedNode) bool) *EnhancedNodes` returns all nodes in the subtree that satisfy the filter function.
- `Len() int` returns the number of nodes in the list wrapper.
- `First() *EnhancedNode` returns the first node in the list wrapper.
- `Last() *EnhancedNode` returns the last node in the list wrapper.
- `At(index int) *EnhancedNode` returns the node at the specified index in the list wrapper, negative indices count from the end.

### Text and Attributes

On `*EnhancedNode`:

- `GetText() (text string, present bool)` returns the text of the first direct text child.
- `GetTexts() []string` returns all texts of the direct text children.
- `ExtractTexts() []string` returns all texts from the full subtree.
- `GetAttr(attr string) (value string, present bool)` returns the value of the specified attribute.
- `RetrieveAttrs(processor func(k, v string) bool)` retrieves all attributes and applies the processor function to each key-value pair. The processor should return `true` to continue processing or `false` to stop.
- `IsTag(tag string) bool` checks if the node's tag name matches the specified tag.
- `HasId(id string) bool` checks if the node has the specified id.
- `HasClasses(classes ...string) bool` checks if the node has all the specified classes.

### Presence and Conversion

On `*EnhancedNode`:

- `Present() bool` checks if the underlying node exists.
- `Unsafe() *UnsafeNode` converts to unchecked access helper.
- `Raw() *html.Node` returns the underlying raw node.
