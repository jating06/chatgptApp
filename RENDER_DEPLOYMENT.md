# Render Deployment Guide

This guide explains how to deploy the MCP server to Render.

## Render Configuration

### 1. Build Command
```bash
go build -tags netgo -ldflags '-s -w' -o app main.go
```

**What it does:**
- `-tags netgo`: Uses Go's native network stack (more portable)
- `-ldflags '-s -w'`: Strips debug information to reduce binary size
- `-o app`: Names the output binary `app`
- `main.go`: Builds from the main.go file

### 2. Pre-Deploy Command (Optional)
```bash
mkdir -p ui
```

**What it does:**
- Ensures the `ui` directory exists for serving HTML widgets
- Creates the directory if it doesn't exist

### 3. Start Command
```bash
./app
```

**What it does:**
- Runs the compiled Go binary
- Server will listen on the port specified by Render's `PORT` environment variable

## Environment Variables

The app automatically uses Render's `PORT` environment variable. No manual configuration needed!

**Optional environment variables you might want to set:**
- None required for basic operation

## Files Required for Deployment

Make sure these files are in your repository:
- ✅ `main.go` - Main server code
- ✅ `go.mod` - Go module definition
- ✅ `go.sum` - Go dependencies checksums
- ✅ `ui/pizza-list.html` - Pizza widget HTML
- ✅ `ui/list-products.html` - Product widget HTML

## Step-by-Step Deployment on Render

### 1. Create New Web Service
- Go to https://render.com/
- Click "New +" → "Web Service"
- Connect your GitHub repository

### 2. Configure Service
**Basic Settings:**
- Name: `mcp-server` (or your preferred name)
- Region: Choose closest to your users
- Branch: `feature/add-pizza-list-widget` or `main`
- Root Directory: Leave empty (or specify if needed)

**Build & Deploy:**
- Runtime: `Go`
- Build Command: `go build -tags netgo -ldflags '-s -w' -o app main.go`
- Pre-Deploy Command: `mkdir -p ui` (optional)
- Start Command: `./app`

**Instance Type:**
- Free (for testing)
- Starter or higher (for production)

### 3. Environment Variables (Optional)
No environment variables are required, but you can add:
- `LOG_LEVEL=info` (if you implement log levels)

### 4. Deploy
- Click "Create Web Service"
- Render will automatically build and deploy your app
- Wait for deployment to complete (usually 2-5 minutes)

## Testing Your Deployment

Once deployed, Render will provide you with a URL like:
```
https://mcp-server-xxxx.onrender.com
```

### Test Endpoints

**1. Health Check:**
```bash
curl https://your-app.onrender.com/health
```

**2. Server Info:**
```bash
curl https://your-app.onrender.com/mcp
```

**3. List Tools:**
```bash
curl -X POST https://your-app.onrender.com/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "tools/list"
  }'
```

**4. Call Pizza List Tool:**
```bash
curl -X POST https://your-app.onrender.com/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 2,
    "method": "tools/call",
    "params": {
      "name": "pizza_list",
      "arguments": {
        "pizzaTopping": "pepperoni"
      }
    }
  }'
```

## Render-Specific Features

### Free Tier Limitations
- Service spins down after 15 minutes of inactivity
- First request after spin-down takes 30-60 seconds
- 750 hours/month of runtime

### Upgrading
For production use, consider:
- **Starter ($7/month)**: No spin-down, better performance
- **Standard+**: Auto-scaling, metrics, alerts

## Port Configuration

The app automatically detects Render's `PORT` environment variable:

```go
func getServerPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}
```

- Render sets `PORT` automatically
- Defaults to `8080` for local development

## Troubleshooting

### Build Fails
- Check that `go.mod` and `go.sum` are committed
- Verify all dependencies are available
- Check build logs for specific errors

### Server Won't Start
- Ensure `./app` has execute permissions
- Check that the binary was created during build
- Review start logs for errors

### 502 Bad Gateway
- Make sure app listens on `0.0.0.0:$PORT`
- Verify the app is binding to the correct port
- Check health endpoint responds

### Images Not Loading
- Verify `ui/` directory is deployed
- Check file paths are correct
- Ensure CSP settings allow external images

## Local Testing with Render-like Environment

Test locally with the PORT environment variable:

```bash
# Build
go build -tags netgo -ldflags '-s -w' -o app main.go

# Run with custom port (simulating Render)
PORT=3000 ./app

# Test
curl http://localhost:3000/health
```

## Continuous Deployment

Render automatically deploys when you push to your connected branch:

```bash
git push origin feature/add-pizza-list-widget
```

Render will:
1. Detect the push
2. Run the build command
3. Run pre-deploy command (if any)
4. Start the new version
5. Route traffic to the new deployment

## Production Checklist

Before going to production:
- [ ] Test all endpoints
- [ ] Verify widgets load correctly
- [ ] Check CSP settings for external resources
- [ ] Set up monitoring/alerts
- [ ] Configure custom domain (optional)
- [ ] Enable HTTPS (automatic on Render)
- [ ] Review and optimize resource usage
- [ ] Consider upgrading from Free tier

## Support

- Render Documentation: https://render.com/docs
- MCP Specification: https://spec.modelcontextprotocol.io/
- Project Issues: Open an issue on GitHub

---

**Your MCP server is now ready to deploy to Render! 🚀**


