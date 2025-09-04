Go through linters.txt and address only the

- tagalign
- lll
- revive (but NOT "avoid meaningless package names")
- staticcheck
- tagalign

warnings, and NOTHING else.

Regularly run `task go:format; task go:lint` to check your progress.
