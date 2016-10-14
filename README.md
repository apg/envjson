# envjson

Run a program in a modified environment, specified via JSON.

## Usage

```
usage: ./envjson [OPTIONS]... JSONfile [COMMAND [ARG]...]

Ensure that the environment meets JSONfile requirements and run COMMAND.

  -i, --ignore-environment  start with an empty environment
      --help     display this help and exit
  -   --stdin    insert read JSON key-value pairs into environment
  -v  --validate-json  validates envjson file

If no COMMAND, print the resulting environment as JSON.
```

By default, `envjson` reads the local environment and merges it with
the JSON file required on the command line. We refer to the provided
JSON file as the "spec". If the spec's requirements are met, e.g., all
required variables are present, it launches the command with an
environment made up of the initial local environment, with defaults
added in from the spec. Therefore, the command:

`envjson requirements.json /usr/local/bin/some-12-factor-app`

ensures that `/usr/local/bin/some-12-factor-app` won't be launched
unless all variables marked "required" in requirements.json are
non-empty.

## Specifying requirements.json

The so called environment "spec" is a single JSON object, with the
following structure:

```json
{
  "ENVIRONMENT_VARIABLE": "default value",
  "ENVIRONMENT_VARIABLE_2": {
    "required": true,
    "inherit": false,
    "value": "<string>",
    "doc": "<string>"
  }
}
```

The meaning of each key is as follows:

* `required`: The command cannot be run if this variable has no value.
* `inherit`: The variable doesn't have a value, and must be inherited
  from the environment.
* `value`: The value that should be used for this variable.
* `doc`: Documentation about the meaning of this environment
  variable. This has no effect on it's value.

## Contributing

Contributions are welcome, and encouraged! Please open an issue before
a Pull Request to avoid duplicated effort, and/or functionality that
will not be merged.

### Thanks

* [Cyril David](https://github.com/cyx) -- for brainstorming.

## Copyright

(c) 2016, Andrew Gwozdziewycz <web@apgwoz.com>

See LICENSE file for more information.
