# GitHub Actions CI/CD Workflows

This directory contains GitHub Actions workflows for automated building, testing, and deployment of the Groceries application.

## Available Workflows

### CI/CD Pipeline (`ci.yml`)

**Purpose**: Automatically build and test both the backend and admin panel on pull requests and pushes to main.

**Triggers**:
- Pull requests to any branch
- Pushes to the `main` branch

**Jobs**:

1. **backend-build**: Builds and tests the Go backend
   - Sets up Go 1.23
   - Downloads and verifies dependencies
   - Generates Swagger documentation
   - Builds the application
   - Runs tests (if any exist)
   - Uploads test coverage

2. **build-summary**: Provides a summary of build results

**Note**: The admin panel has its own separate CI workflow in the [groceries-admin](https://github.com/superbkibbles/groceries-admin) repository.

## Setup Requirements

### For CI (Current Implementation)

No special setup required! The workflow will automatically run when:
- You create a pull request
- You push to the main branch

### For Admin Panel

The admin panel has its own CI workflow in a separate repository. See the [groceries-admin](https://github.com/superbkibbles/groceries-admin) repository for its CI configuration.

## Future Deployment Setup

To enable automatic deployment to your VPS, you'll need to:

### 1. Generate SSH Key Pair

On your local machine:
```bash
ssh-keygen -t ed25519 -C "github-actions" -f ~/.ssh/github-actions
```

### 2. Add Public Key to VPS

Copy the public key to your VPS:
```bash
ssh-copy-id -i ~/.ssh/github-actions.pub user@your-vps-ip
```

### 3. Add GitHub Secrets

Go to your repository → Settings → Secrets and variables → Actions, and add:

- `SSH_PRIVATE_KEY`: Contents of `~/.ssh/github-actions` (private key)
- `VPS_HOST`: Your VPS IP address (e.g., `192.168.1.100`)
- `VPS_USER`: SSH username (e.g., `ubuntu`, `root`, etc.)
- `VPS_PORT`: SSH port (default: `22`)

### 4. Extend Workflow for Deployment

Add a deployment job that:
- Builds Docker images
- Pushes to registry or directly to VPS
- SSH into VPS and run docker-compose

## Monitoring Builds

### View Build Status

1. Go to your repository on GitHub
2. Click on the "Actions" tab
3. See all workflow runs and their status

### Add Status Badge to README

Add this to your `README.md`:

```markdown
![CI/CD Pipeline](https://github.com/YOUR_USERNAME/groceries-backend/workflows/CI%2FCD%20Pipeline/badge.svg)
```

Replace `YOUR_USERNAME` with your GitHub username.

## Troubleshooting

### Go Build Fails

**Issue**: Dependencies cannot be downloaded

**Solutions**:
- Verify `go.mod` and `go.sum` are committed
- Check for any private module dependencies
- Ensure Go version matches (currently set to 1.23)

## Local Testing

Test your workflow locally before pushing using [act](https://github.com/nektos/act):

```bash
# Install act
brew install act  # macOS
# or
curl https://raw.githubusercontent.com/nektos/act/master/install.sh | sudo bash

# Run workflow
act pull_request
```

## Customization

### Change Go Version

Edit line 23 in `ci.yml`:
```yaml
go-version: '1.23'  # Change to your preferred version
```

### Add Environment Variables for Build

Add to the build step:
```yaml
- name: Build application
  run: go build -v -o build/groceries-api ./main.go
  env:
    CGO_ENABLED: "0"
    GOOS: linux
```

## Best Practices

1. **Always test locally first**: Before pushing, ensure your code builds locally
2. **Keep secrets secure**: Never commit secrets or credentials
3. **Review failed builds**: Check logs to understand why builds fail
4. **Update dependencies**: Keep actions up to date (currently using v4/v5)
5. **Monitor build times**: Optimize if builds take too long

## Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Go Setup Action](https://github.com/actions/setup-go)
- [Node Setup Action](https://github.com/actions/setup-node)
- [Deploying with GitHub Actions](https://docs.github.com/en/actions/deployment)

