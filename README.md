# pincode-api

A small Go web project (templates included) — repository contains a simple Go server with templates under `templates/` and the main application in `main.go`.

## Contents

- `main.go` — application entrypoint
- `templates/index.html` — HTML template used by the server

## Prerequisites

- Go 1.20+ (Go 1.21 recommended)
- Git (for version control)

## Quick start (run locally)

Open a Windows Command Prompt (cmd.exe) and run from the project root:

```bat
cd "c:\Users\K YADU VAMSI KRISHNA\Documents\go projects\Test GO"

# Run directly with `go run`
go run main.go

# Or build a binary and run it
go build -o pincode-api
.pincode-api.exe  # or .\pincode-api.exe
```

The app will typically listen on the port configured in `main.go` (check source for exact host:port). Visit `http://localhost:PORT` in your browser (replace PORT with the configured port).

## Git / repository notes

If you haven't set your Git identity yet, configure it (recommended globally once):

```bat
git config --global user.email "your-email@example.com"
git config --global user.name "Your Name"
```

If you only want to set the identity for this repo, omit `--global`.

Common Git commands for this repo:

```bat
# see status
git status

# stage everything and commit
git add .
git commit -m "pincode"

# show last few commits
git log --oneline -5

# rename local branch from master to main (optional)
git branch -M main

# push to remote (set upstream)
git push -u origin main
```

Note: you earlier saw the repository reporting `On branch master` while platform metadata showed `main` — that's purely a local branch-name difference. Renaming to `main` above will align local branch name with the conventional `main` branch if you prefer.

## Project structure

```
Test GO/
├─ go.mod
├─ main.go
├─ README.md
└─ templates/
   └─ index.html
```

## Troubleshooting

- "'git' is not recognized" — install Git for Windows and select the option to add Git to PATH, then open a new cmd window.
- After changing PATH or installing Git, restart your terminal so the change takes effect.
- If `git commit` complains about unknown identity, run the `git config` commands above.

## Contributing

Contributions are welcome. Please open issues or PRs with small, focused changes.

## License

Add a license file if you plan to make this project public (MIT, Apache-2.0, etc.).
