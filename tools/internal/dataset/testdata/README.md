# Validator fixtures

Every species here (`Testus alpha`, `Testus beta`, …) is fictional, and so are its values. These files only exercise the validator; never copy them into `fish/`.

- `valid/`: a dataset that must pass.
- `invalid/<case>/`: one dataset per failure case, used by `validate_test.go`.
- `plants/valid/` and `plants/invalid/<case>/`: the same for the plant schema (`Testus gamma`, `Testus delta`), used by `plants_test.go`.
