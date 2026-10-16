# marlin-dataset

An open dataset of **aquarium fish**.

## What's inside

```
fish/<genus>/<id>.json       # one species per file
schema/fish.schema.json      # format rules
template.json                # template for new species
tools/                       # Go CLI for validation
```

Example: [`fish/paracheirodon/paracheirodon-innesi.json`](fish/paracheirodon/paracheirodon-innesi.json).

- Metric units: cm, liters, °C, hardness in °dH.
- Ranges as `{ min, max }`.
- `null` means "not known with certainty". It is never filled with an approximate value.
- `status: draft` = not reviewed; `status: reviewed` = verified by a person and backed by sources.
- Common names are provided in Spanish (`es`, required) and English (`en`).

## Using the data

Read the JSON files directly and load them into whatever tool or database you prefer. Before showing a value to end users, filter by `status: reviewed`.

### License and attribution

You can use the data and code for anything, including commercial purposes. The [MIT license](LICENSE) only asks you to keep the copyright notice. A credit line is appreciated:

> Data from marlin-dataset (https://github.com/srteerra/marlin-dataset), MIT License.

## Validating

Requires Go (version in [`tools/go.mod`](tools/go.mod)).

```bash
cd tools
go test ./...              # tests
go run . validate ../fish  # validate the whole dataset
```

All errors are reported together, with file, field and reason.

## Contributing

Welcome! Read [CONTRIBUTING.md](CONTRIBUTING.md). If you don't use Git, open an issue with the **Species** template.

This project follows a [Code of Conduct](CODE_OF_CONDUCT.md).

## License

Everything in this repository (data, schema, documentation and code) is licensed under the [MIT License](LICENSE).
