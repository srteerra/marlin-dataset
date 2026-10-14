# Contributing

Thanks for helping! No prior experience needed: if you know about fish, you can already contribute.

## Adding or fixing a species

1. **Copy the template.** Copy [`template.json`](template.json) to `fish/<genus>/<id>.json`.
   - The `id` is the scientific name in lowercase, with hyphens and no accents.
   - The folder is the first part of the `id`.
   - Example: *Paracheirodon innesi* → `fish/paracheirodon/paracheirodon-innesi.json`.
   - Keep the `"$schema"` line as is: editors like VS Code use it to autocomplete fields, show what each one means and flag mistakes as you type.
2. **Fill in the data** using your sources and your own words.
   - Metric units: cm, liters, °C, °dH.
   - If you are not sure about a value, leave it as `null`. Made-up data can kill real fish.
   - Ranges look like `{ "min": 22, "max": 26 }`.
   - Write `notes` and `reason` in English. Common names go in `common_names.es` (Spanish, at least one) and `common_names.en` (English).
   - Put your sources in `sources`.
   - Keep `"status": "draft"`. A person who reviews the data will change it to `"reviewed"`.
3. **Validate.** If you have Go installed:
   ```bash
   cd tools
   go run . validate ../fish
   ```
   If not, no problem: CI validates your PR automatically.
4. **Open a pull request.** One species per PR, so it's easy to review.

### Allowed values

| Field | Values |
|---|---|
| `water_type` | `freshwater`, `brackish`, `marine` |
| `temperament` | `peaceful`, `semi_aggressive`, `aggressive` |
| `swim_zone` | `top`, `middle`, `bottom`, `all` |
| `diet` | `herbivore`, `carnivore`, `omnivore` |
| `care_level` | `beginner`, `intermediate`, `advanced` |
| `schooling_min` | minimum group size (2 or more), `null` if it does not live in groups |
| `incompatible_with` | only what cannot be deduced from the other fields, e.g. `[{ "id": "pterophyllum-scalare", "reason": "Adult angelfish may eat it." }]` |

Every field is described in [`schema/fish.schema.json`](schema/fish.schema.json).

### Images

Images are **links, never files**. Do not add image files to the repository.

- For now only [Pixabay](https://pixabay.com) is accepted.
- Use the image **page** URL (e.g. `https://pixabay.com/photos/neon-tetra-1234567/`), not a direct file link.
- Make sure the photo really shows this species: many photos online are misidentified.
- Avoid photos with visible brands, logos or recognizable people.

```json
"images": [
  { "url": "https://pixabay.com/photos/neon-tetra-1234567/", "author": "Author name", "license": "pixabay" }
]
```

### Don't use Git?

Open an issue with the **Species** template: write what you know and your sources, and someone will turn it into a file.

## Important rules

- **Do not copy text** from websites, books or databases. Facts (a temperature range, a size) can be used; the wording has to be yours.
- Do not extract tables from databases with non-commercial or share-alike licenses (for example FishBase). You can consult them to verify a fact, nothing more.
- No image files: only links to images (see [Images](#images)).
- Need a new field? Propose it in an issue first; schema changes go in a separate PR.

## Commit messages

We use [Conventional Commits](https://www.conventionalcommits.org/):

```
feat(fish): add paracheirodon-innesi
fix(fish): correct ph for betta-splendens
feat(tools): add hardness check
docs: ...
ci: ...
```

## License of your contributions

By contributing, you agree that your contribution is published under the same license as the repository: the **MIT License**, for data, documentation and code alike. This follows the "inbound = outbound" rule in [GitHub's Terms of Service](https://docs.github.com/en/site-policy/github-terms/github-terms-of-service#6-contributions-under-repository-license). Nothing else to sign.

## Conduct

This project follows a [Code of Conduct](CODE_OF_CONDUCT.md). To report a problem, email srterradev@gmail.com.
