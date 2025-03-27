# fngrep
## Installation
Clone this repo and run `go install`, make sure your golang is properly installed. Check if you can run binary under $GOPATH/bin.

## What is fngrep
grep function body from random file

## Usage
```
Usage:
  fngrep [OPTIONS] Filename

Application Options:
  -p, --prefix=   Function prefix to grep
  -r, --regexp=   Function line regex to grep

Help Options:
  -h, --help      Show this help message
```

## Example
```
> fngrep -p "type ActivateProductResponse {" somefile.graphql
type ActivateProductResponse {
  error: [json]
  id_list: [String]
  status: String!
}

> fngrep -r '\s+warehouse_location\(' somefile.graphql
  warehouse_location(
    distinct_on: [warehouse_location_select_column!]
    limit: Int
    offset: Int
    order_by: [warehouse_location_order_by!]
    where: warehouse_location_bool_exp
  ): [warehouse_location!]!
```
