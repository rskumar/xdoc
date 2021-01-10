# xdoc

Structured Document Content Schema implementation

## Implemented Elements


### Page `<page />`

### Title `<title />`

### Paragraph `<para />`

## Schema

[XML Schema](extra/schema.ref.xml)

[JSON Schema](extra/schema.ref.json)

## Implementing a new Element

- A new element should satisfy interface `Element`
- It's mapping entry should be in `schema` map in `schema.go`.
- In `tree.go`, add its entry in `switch` inside marshal and unmarshal.
- Any data directly related to element should be its `struct` attribute. Example -
  ```go
  type Text struct {
    Text string
    Bold bool
  } 
  ```
- If element supports children, it should embed `Children` with proper json and xml tags for their correct marshalling
  and unmarshalling. Example -
  ```go
  type Para struct {
    Children `json:"children" xml:",any"`
  }
  ```

- Custom element specific Marshal/Unmarshal - To customize the JSON and XML marshalling and unmarshalling, in new
  element implement/satisfy any or all of -
  - `json.Marshaller`
  - `json.Unmarshaller`
  - `xml.Marshaller`
  - `xml.Unmarshaller`

  If element doesn't satisfy any, the default behaviour is implied, which is purely based on struct tags in struct
  fields (`json` or `xml`)
  

