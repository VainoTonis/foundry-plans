# foundry-plans

A minimal Go CLI for managing plans via Foundry HTTP API.

## Build

```bash
make build
```

## Install

```bash
make install
```

Installs `foundry-plans` to `~/.local/bin/foundry-plans`.

## Usage

```bash
./foundry-plans [flags] <command>
```

### Global Flags

- `--url`: Foundry API URL (default: `http://localhost:8080`)

### Commands

#### List Plans (read-only)

```bash
./foundry-plans plans list
```

#### Create Plan (JSON-only)

Create a new plan. All fields must come from JSON stdin. No positional arguments.

Required JSON fields: `repo_name`, `title`
Optional JSON fields: `summary`, `steps` (array of strings)

```bash
echo '{
  "repo_name": "my-repo",
  "title": "My Plan",
  "summary": "Plan description",
  "steps": ["step1", "step2", "step3"]
}' | ./foundry-plans plans create
```

#### Get Plan (read-only)

```bash
./foundry-plans plans get 123
```

#### Update Plan (JSON-only)

Update a plan. All fields including the plan `id` must come from JSON stdin. No positional arguments.

Required JSON field: `id`
Optional JSON fields: `status`, `title`, `summary`

```bash
echo '{
  "id": 123,
  "status": "in_progress",
  "title": "Updated Title"
}' | ./foundry-plans plans update
```

#### Update Step (JSON-only)

Update a step. All fields including `plan_id` and either `step_id` or `position` must come from JSON stdin. No positional arguments.

Required JSON fields: `plan_id`, and either `step_id` (number) or `position` (number)
Optional JSON fields: `status`, `text`

You can identify a step by `step_id` (the unique step ID):

```bash
echo '{
  "plan_id": 123,
  "step_id": 456,
  "status": "completed",
  "text": "Step completed successfully"
}' | ./foundry-plans plans update-step
```

Or by `position` (the step's position in the plan, 1-based):

```bash
echo '{
  "plan_id": 123,
  "position": 1,
  "status": "completed",
  "text": "First step completed"
}' | ./foundry-plans plans update-step
```

## Examples

```bash
# List all plans
./foundry-plans plans list

# Create a plan with steps
echo '{
  "repo_name": "example",
  "title": "Setup",
  "summary": "Initial setup plan",
  "steps": ["install", "configure", "deploy"]
}' | ./foundry-plans plans create

# Get plan details
./foundry-plans plans get 123

# Update plan status and title
echo '{
  "id": 123,
  "status": "in_progress",
  "title": "Setup in Progress"
}' | ./foundry-plans plans update

# Update a step by step_id
echo '{
  "plan_id": 123,
  "step_id": 456,
  "status": "completed",
  "text": "Installation complete"
}' | ./foundry-plans plans update-step

# Update a step by position
echo '{
  "plan_id": 123,
  "position": 1,
  "status": "completed",
  "text": "First step done"
}' | ./foundry-plans plans update-step

# Use custom API URL
./foundry-plans --url http://api.example.com plans list
```

## Output

All responses are formatted as pretty-printed JSON.

## API Endpoints

The CLI uses the following Foundry API endpoints:

- `GET /api/plans` - List all plans
- `POST /api/plans` - Create a plan
- `GET /api/plans/{id}` - Get plan details
- `GET /api/plans/{id}/steps` - Get plan steps
- `POST /api/plans/{id}/steps` - Create a step
- `PATCH /api/plans/{id}` - Update plan (status, title, summary)
- `PATCH /api/plans/{id}/steps/{step_id}` - Update step (status, text)
