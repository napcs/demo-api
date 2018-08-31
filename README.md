# demo-api

A simple command-line JSON server. If you need to stub out an API locally, this is the tool you need. It's lightweight, works on Mac, Windows, and Linux, and requires no setup. 

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

Then run `demo-api`. 

Visit `http://localhost:8080` in the browser to see the file served as JSON.

In addition, since the JSON file is structured like a typical REST API data structure, the following will work:

* `GET http://localhost:8080/notes` - retrive all "notes"
* `GET http://localhost:8080/notes/1` - retrive the note with the `id`   of `1`
* `POST http://localhost:8080/notes/` - Create a new note. This modifies the `data.json` file.
* `DELETE http://localhost:8080/notes/1` - Delete the note with the `id`   of `1`. This modifies the `data.json` file.

## Installation

To install, download the latest release to your system and copy the executable to a location on your path. Then launch it in a directory containing `data.json`.


## Roadmap

* `PUT/PATCH` support
* tests
* refactoring
* help switch
* version notice

## Contributing

Please contribute. Fork, create a PR.

## License

Apache 2. See LICENSE file.
