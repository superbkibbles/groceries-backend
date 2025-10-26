# GitHub Actions Setup Guide

This guide will help you set up GitHub Actions CI/CD for your Groceries project.

## ✅ What's Already Done

The following files have been created and configured:

1. **`.github/workflows/ci.yml`** - Main CI/CD workflow
2. **`.github/workflows/README.md`** - Detailed documentation
3. **`README.md`** - Updated with CI badge and correct Go version

## 🚀 Getting Started

### Step 1: Push to GitHub

If you haven't already, push these changes to your GitHub repository:

```bash
cd /Users/superbkibbles/Documents/projects/groceries/groceries-backend

# Add the new workflow files
git add .github/
git add README.md

# Commit the changes
git commit -m "Add GitHub Actions CI/CD workflow"

# Push to GitHub
git push origin main
```

### Step 2: Verify Admin Repository

The workflow assumes your admin panel is in a separate repository. Verify:

1. The repository exists at: `https://github.com/YOUR_USERNAME/groceries-admin`
2. If the name is different, update line 83 in `.github/workflows/ci.yml`:
   ```yaml
   repository: YOUR_USERNAME/YOUR_ADMIN_REPO_NAME
   ```

### Step 3: Test the Workflow

Create a test pull request to trigger the workflow:

```bash
# Create a new branch
git checkout -b test-ci

# Make a small change
echo "# Test CI" >> TEST.md
git add TEST.md
git commit -m "Test CI workflow"

# Push and create PR
git push origin test-ci
```

Then go to GitHub and create a pull request. The workflow will automatically run!

### Step 4: Check Workflow Results

1. Go to your repository on GitHub
2. Click the **"Actions"** tab
3. You'll see the "CI/CD Pipeline" workflow running
4. Click on it to see detailed logs

## 📊 What the Workflow Does

### On Every Pull Request:
- ✅ Builds the Go backend
- ✅ Generates Swagger documentation
- ✅ Runs backend tests
- ✅ Builds the Next.js admin panel
- ✅ Runs admin panel linting and tests

### On Push to Main:
- Same as above
- Ready to add deployment steps

## 🔧 Common Issues & Solutions

### Issue: Admin Repository Not Found

**Error**: `Repository not found` or `403 Forbidden`

**Solution**:
1. Verify the admin repository exists
2. Check the repository name in `.github/workflows/ci.yml` line 83
3. For private repos, ensure GitHub Actions has access

### Issue: Build Fails on Go Dependencies

**Error**: `go: github.com/some/package: module not found`

**Solution**:
1. Ensure `go.mod` and `go.sum` are committed
2. Run `go mod tidy` locally and commit changes
3. Check for private packages that need authentication

### Issue: Next.js Build Fails

**Error**: Build errors in Next.js

**Solution**:
1. Test build locally: `cd ../groceries-admin && npm run build`
2. Fix any build errors
3. Commit and push fixes

## 🎯 Next Steps: Deployment

Once CI is working, you can add automatic deployment to your VPS:

### 1. Generate SSH Keys

On your local machine:
```bash
ssh-keygen -t ed25519 -C "github-actions-deploy" -f ~/.ssh/github-actions
```

### 2. Add Public Key to VPS

```bash
# Copy public key to clipboard
cat ~/.ssh/github-actions.pub

# SSH into your VPS
ssh your_user@your_vps_ip

# Add the key to authorized_keys
echo "PASTE_PUBLIC_KEY_HERE" >> ~/.ssh/authorized_keys

# Set correct permissions
chmod 600 ~/.ssh/authorized_keys
chmod 700 ~/.ssh
```

### 3. Add GitHub Secrets

Go to: Repository → Settings → Secrets and variables → Actions

Add these secrets:
- **`SSH_PRIVATE_KEY`**: Contents of `~/.ssh/github-actions` (the PRIVATE key)
- **`VPS_HOST`**: Your VPS IP (e.g., `192.168.1.100`)
- **`VPS_USER`**: SSH username (e.g., `ubuntu`)
- **`VPS_PORT`**: SSH port (default: `22`)

### 4. Add Deployment Job to Workflow

Create a new file `.github/workflows/deploy.yml`:

```yaml
name: Deploy to VPS

on:
  push:
    branches:
      - main
  workflow_dispatch:  # Allow manual trigger

jobs:
  deploy:
    name: Deploy to Production
    runs-on: ubuntu-latest
    
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
      
      - name: Deploy to VPS
        uses: appleboy/ssh-action@master
        with:
          host: ${{ secrets.VPS_HOST }}
          username: ${{ secrets.VPS_USER }}
          key: ${{ secrets.SSH_PRIVATE_KEY }}
          port: ${{ secrets.VPS_PORT }}
          script: |
            cd /path/to/your/app
            git pull origin main
            docker-compose down
            docker-compose up -d --build
            docker-compose ps
```

## 📚 Additional Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Workflow Syntax](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions)
- [Self-hosted Runners](https://docs.github.com/en/actions/hosting-your-own-runners)

## 🆘 Getting Help

If you encounter issues:
1. Check the Actions tab for detailed error logs
2. Review the [workflows README](.github/workflows/README.md)
3. Search GitHub Actions documentation
4. Check if the issue is related to your code or the workflow

## ✨ Success Indicators

You'll know everything is working when:
- ✅ The CI badge in README.md shows "passing"
- ✅ Pull requests show green checkmarks
- ✅ Both backend and admin builds complete successfully
- ✅ No red X marks in the Actions tab

Happy coding! 🚀

