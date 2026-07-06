# foundry-plans

A minimal Go CLI for managing plans via Foundry HTTP API.

## Build

```bash
make build
```

## Usage

```bash
./foundry-plans [flags] <command>
```

### Global Flags

- `--url`: Foundry API URL (default: `http://localhost:8080`)

### Commands

#### List Plans

```bash
./foundry-plans plans list
```

#### Create Plan

```bash
./foundry-plans plans create \
  --repo-name my-repo \
  --title "My Plan" \
  --summary "Plan description" \
  --steps "step1,step2,step3"
```

#### Get Plan

```bash
./foundry-plans plans get 123
```

#### Update Plan Status

```bash
./foundry-plans plans update-status 123 in_progress
```

#### Update Step

```bash
./foundry-plans plans update-step 123 456 completed
```

## Examples

```bash
# List all plans
./foundry-plans plans list

# Create a plan with steps
./foundry-plans plans create \
  --repo-name example \
  --title "Setup" \
  --summary "Initial setup plan" \
  --steps "install,configure,deploy"

# Get plan details
./foundry-plans plans get 123

# Update plan status
./foundry-plans plans update-status 123 in_progress

# Update a step
./foundry-plans plans update-step 123 456 completed "Step completed successfully"

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
- `PATCH /api/plans/{id}` - Update plan status
- `PATCH /api/plans/{id}/steps/{step_id}` - Update step status/text
