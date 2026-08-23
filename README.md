# foundry-plans

> Deprecated. Plan commands now live in Foundry's supported `foundry` CLI at
> `cmd/foundry`. This repository remains available for history only.

## Migration

Build and install the replacement from the Foundry repository:

```sh
cd /path/to/foundry
go build -o ~/.local/bin/foundry ./cmd/foundry
```

The default API URL remains `http://localhost:8080`; override it with
`foundry --url <url>`.

Resolve repository IDs:

```sh
foundry repositories list
```

Create a plan. `repository_ids` is ordered and non-empty, and step positions are
zero-based:

```sh
printf '%s\n' '{
  "repository_ids": [1],
  "title": "My plan",
  "summary": "Goal, scope, non-goals, decisions, and risks",
  "steps": [
    {"text": "First step", "parallel_group": 1},
    {"text": "Second step", "parallel_group": 1}
  ]
}' | foundry plans create
```

Read and update plans:

```sh
foundry plans list
foundry plans get <plan_id>
printf '%s\n' '{"status":"running"}' | foundry plans update <plan_id>
printf '%s\n' '{"plan_id":<plan_id>,"position":0,"status":"done"}' | foundry plans update-step
```

Run a plan:

```sh
foundry plans run <plan_id>
```

Do not add new features here. Extend `cmd/foundry` and `internal/apiclient` in
the Foundry repository instead.
