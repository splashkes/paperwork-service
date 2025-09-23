# Deployment Guide - Art Battle Paperwork Service

## Prerequisites

1. **DigitalOcean CLI (doctl)**
   ```bash
   # Install doctl
   curl -sL https://github.com/digitalocean/doctl/releases/download/v1.100.0/doctl-1.100.0-linux-amd64.tar.gz | tar -xzv
   sudo mv doctl /usr/local/bin

   # Authenticate
   doctl auth init
   ```

2. **Supabase Edge Function**
   The service depends on the `paperwork-data` edge function. Deploy it first:
   ```bash
   cd /root/vote_app/vote26
   supabase functions deploy paperwork-data
   ```

## Deployment Steps

### 1. Initial Setup

1. **Build locally to test**:
   ```bash
   cd /root/vote_app/paperwork-service
   make deps
   make build
   ```

2. **Test locally** (optional):
   ```bash
   # Create .env file with your Supabase credentials
   cp .env.example .env
   # Edit .env with your actual values

   # Run locally
   make run

   # Test health check
   curl http://localhost:8080/api/v1/health

   # Test PDF generation (replace AB2995 with actual event EID)
   curl "http://localhost:8080/api/v1/event-pdf/AB2995" -o test.pdf
   ```

### 2. Deploy to DigitalOcean

1. **Run deployment script**:
   ```bash
   ./deploy.sh
   ```

2. **Set environment variables** in DigitalOcean Apps dashboard:
   - `SUPABASE_URL`: Your Supabase project URL
   - `SUPABASE_KEY`: Your Supabase anon key
   - `ENVIRONMENT`: `production`

3. **Monitor deployment**:
   - Check the DigitalOcean Apps dashboard
   - Watch for build completion and successful health checks

### 3. Post-Deployment

1. **Test the deployed service**:
   ```bash
   # Health check
   curl https://your-app-url.ondigitalocean.app/api/v1/health

   # PDF generation
   curl "https://your-app-url.ondigitalocean.app/api/v1/event-pdf/AB2995" -o test.pdf
   ```

2. **Configure custom domain** (optional):
   - Add domain in DigitalOcean Apps dashboard
   - Update DNS records

## Environment Variables

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `SUPABASE_URL` | Supabase project URL | ✅ | - |
| `SUPABASE_KEY` | Supabase anon key | ✅ | - |
| `PORT` | Server port | ❌ | `8080` |
| `ENVIRONMENT` | Environment (development/production) | ❌ | `development` |
| `TEMPLATES_PATH` | Path to templates directory | ❌ | `./templates` |

## API Endpoints

- `GET /api/v1/health` - Health check
- `GET /api/v1/event-pdf/{eid}` - Generate PDF for event (e.g., `/api/v1/event-pdf/AB2995`)

## Monitoring

### Health Checks
The service includes automatic health checks at `/api/v1/health`. DigitalOcean will monitor this endpoint.

### Logs
View application logs in the DigitalOcean Apps dashboard or via doctl:
```bash
doctl apps logs <app-id> --follow
```

### Metrics
Monitor CPU, memory, and request metrics in the DigitalOcean Apps dashboard.

## Troubleshooting

### Common Issues

1. **Edge function not found**:
   - Ensure `paperwork-data` edge function is deployed in vote26 project
   - Check Supabase function logs

2. **PDF generation fails**:
   - Check template files are included in Docker image
   - Verify font files are accessible
   - Check application logs for specific errors

3. **Event not found**:
   - Verify event EID exists in Supabase
   - Check event is enabled (`enabled = true`)
   - Ensure event has artists assigned

4. **Build failures**:
   - Run `make deps` to ensure dependencies are up to date
   - Check Go version compatibility
   - Verify all required files are included

### Debug Commands

```bash
# Check app status
doctl apps get <app-id>

# View logs
doctl apps logs <app-id> --follow

# Force rebuild
doctl apps create-deployment <app-id> --force-rebuild

# Get app URL
doctl apps get <app-id> --format LiveURL --no-header
```

## Scaling

The service is configured with minimal resources (`basic-xxs`). For high-traffic scenarios:

1. **Increase instance size** in `.do/app.yaml`:
   ```yaml
   instance_size_slug: basic-xs  # or basic-s
   ```

2. **Enable auto-scaling**:
   ```yaml
   instance_count: 1
   autoscaling:
     min_instance_count: 1
     max_instance_count: 3
   ```

3. **Add performance monitoring**:
   - Monitor response times
   - Watch memory usage during PDF generation
   - Set up alerts for high error rates

## Security

- Service uses CORS middleware for cross-origin requests
- No authentication required (public endpoint)
- Supabase access is read-only via edge function
- All data fetched from public Supabase endpoints

## Maintenance

### Updates
1. Make code changes
2. Test locally
3. Commit to repository
4. Run `./deploy.sh` to deploy

### Template Updates
To update PDF templates (fonts, backgrounds):
1. Replace files in `templates/` directory
2. Rebuild and redeploy service

### Edge Function Updates
If the Supabase edge function changes:
1. Deploy updated function to Supabase
2. No service restart required