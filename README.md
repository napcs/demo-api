# demo-api

A simple command-line JSON server that also emulates a REST API. If you need to stub out an API locally, this is the tool you need. It's lightweight, works on Mac, Windows, and Linux, and requires no setup other than a JSON data file. 

## Usage

Create a JSON file called `data.json`:

```json
{
  "notes" : [
    {
      "id" : 1,
      "title" : "Hello World"
    },
    {
      "id" : 2,
      "title" : "The second note"
    }
  ]
}
```

Then run `demo-api`.  The server starts up. 

Visit `http://localhost:8080` in the browser to serve the entire file as a JSON response.

In addition, since the JSON file is structured like a typical REST API data structure, the following endpoints will work:

* `GET http://localhost:8080/notes` - retrive all "notes"
* `GET http://localhost:8080/notes/1` - retrive the note with the `id`   of `1`
* `POST http://localhost:8080/notes/` with JSON payload - Create a new note. This modifies the `data.json` file.
* `PUT/PATCH http://localhost:8080/notes/1` with JSON payload - update the note with the `id`   of `1`. This modifies the `data.json` file.
* `DELETE http://localhost:8080/notes/1` - Delete the note with the `id`   of `1`. This modifies the `data.json` file.


Anything that doesn't match returns a `404` status code.

## Examples with `curl`

Given the following data file:


```json
{
  "notes" : [
    {
      "id" : 1,
      "title" : "Hello World"
    },
    {
      "id" : 2,
      "title" : "The second note"
    }
  ]
}
```

To get everything:

```
curl -i localhost:8080/
```

To get the `notes` node:


```
curl -i localhost:8080/notes
```

To get the `notes/1` node:


```
curl -i localhost:8080/notes/1
```

To add a new note:

```
curl -i -X POST http://localhost:8080/notes \
-H "Content-type: application/json" \
-d '{"title": "This is another note"}'
```

To update the contents of the first note:

```
curl -i -X PUT http://localhost:3000/notes/1  \
-H "Content-type: application/json" \
-d '{"title": "This is the third note"}'
```

To delete the third note:

```
curl -i -X DELETE localhost:3000/notes/3
```

If you use a different JSON file, your paths will be different.

### Advanced Usage

To specify a different port, use the `-p` option:

```
demo-api -p 4000
```

To specify a different filename, in case you don't like `data.json` as the default, use `-f` and specify the file:


```
demo-api -f notes.json
```

To view the version, use `-v`:

```
demo-api -v
```

## Installation

To install, download the latest release to your system and copy the executable to a location on your path. Then launch it in a directory containing `data.json`.


## Roadmap

* Refactoring. This code is a mess.
* A "no persist" mode - changes are accepted but not saved to the JSON file.

## Contributing

Please contribute.

Clone the repository and then download the dependencies:

```
$ go get github.com/Jeffail/gabs
$ go get github.com/gin-gonic/gin
$ go get github.com/codegangsta/gin
```

Run development version:

```
$ gin --appPort 8080 go run app.go
```

The server is now listening on `localhost:3000` and will reload on code change.

Make changes, run the tests, create a PR. 

## History

* 2018-09-03 - v0.3.0
  * Refactoring to make testing possible
  * Adds test suite
  * Pretty print
  * Supports `PUT` and `PATCH`
  * Supports `-v` option to show version
  * Supports `-f` option to specify the data file
  * Supports `-p` option to specify the port
  * Fix bug where querying non-existant ID still returned a 200 status code instead of 404
* 2018-08-30 - v0.1.0
  * initial release

## License

Apache 2. See LICENSE file.
