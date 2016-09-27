# envjson

Run a program in a modified environment, specified via JSON.

## Usage

`usage: envjson [--stdin | --help] JSONfile <cmd> <cmd args>`

By default, `envjson` reads the local environment and merges it with
the JSON file required on the command line. We refer to the provided
JSON file as the "spec". If the spec's requirements are met, e.g., all
required variables are present, it launches the command with an
environment made up of the initial local environment, with defaults
added in from the spec.

`envjson requirements.json /usr/local/bin/some-12-factor-app`

Additionally, the initial environment can be augmented by reading JSON
from STDIN, in which case, the key-value pairs read will overwrite the
local environment before checking requirements, and merging defaults.

`envjson --stdin requirements.json /usr/local/bin/some-12-factor-app`

## Specifying requirements.json

The so called environment "spec" is a single JSON object, with the
following structure:

```json
{
  "ENVIRONMENT_VARIABLE": "default value",
  "ENVIRONMENT_VARIABLE_2": {
    "required": <bool>,
    "inherit": <bool>,
    "default": <string>,
    "doc": <string>
  }
}
```

The meaning of each key is as follows:

* `required`: The command cannot be run if this variable has no value.
* `inherit`: The variable doesn't have a default value, and must be
  inherited from the local environment
* `default`: The default value for this variable.
* `doc`: Documentation about the meaning of this environment
  variable. This has no effect on it's value.

## Contributing

Contributions are welcome, and encouraged! Please open an issue before
a Pull Request to avoid duplicated effort, and/or functionality that
will not be merged.

## Copyright

(c) 2016, Andrew Gwozdziewycz <web@apgwoz.com>

See LICENSE file for more information.
