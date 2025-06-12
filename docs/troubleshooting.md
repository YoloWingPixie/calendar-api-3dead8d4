# Troubleshooting Guide

## Common Issues and Solutions

### Database Connection Issues

#### Symptoms
- API fails to start
- Health check shows database as disconnected
- 500 errors when accessing endpoints

#### Solutions
1. Check database credentials in environment variables
2. Verify database is running: `docker ps | grep postgres`
3. Check database logs: `docker logs calendar-api-postgres`
4. Verify network connectivity: `docker network inspect calendar-api_default`

### API Authentication Issues

#### Symptoms
- 401 Unauthorized errors
- API key not being accepted
- Missing X-API-Key header

#### Solutions
1. Verify API key is correctly set in environment
2. Check API key format and length
3. Ensure X-API-Key header is being sent
4. Check API logs for authentication errors

### Container Issues

#### Symptoms
- Container fails to start
- Container crashes repeatedly
- Port conflicts

#### Solutions
1. Check container logs: `docker logs calendar-api-server`
2. Verify port availability: `netstat -tulpn | grep 8000`
3. Check container health: `docker ps -a`
4. Verify environment variables: `docker exec calendar-api-server env`

### Development Environment Issues

#### Symptoms
- Build failures
- Test failures
- Dependency issues

#### Solutions
1. Clean and rebuild: `task clean && task build`
2. Update dependencies: `task mod`
3. Check Go version: `go version`
4. Verify Docker version: `docker version`

## Performance Issues

### High Response Times

#### Symptoms
- Slow API responses
- Timeout errors
- High database load

#### Solutions
1. Check database connection pool settings
2. Monitor database performance
3. Review query execution plans
4. Check for missing indexes

### Memory Issues

#### Symptoms
- Container OOM errors
- High memory usage
- Slow response times

#### Solutions
1. Adjust container memory limits
2. Monitor memory usage: `docker stats`
3. Review application memory profile
4. Check for memory leaks

## Deployment Issues

### CI/CD Pipeline Failures

#### Symptoms
- Pipeline jobs failing
- Deployment errors
- Test failures

#### Solutions
1. Check GitHub Actions logs
2. Verify secrets and environment variables
3. Review test output
4. Check infrastructure state

### Infrastructure Issues

#### Symptoms
- Terraform apply failures
- Resource creation errors
- Configuration issues

#### Solutions
1. Check Terraform state
2. Verify AWS credentials
3. Review resource limits
4. Check network configuration

## Getting Help

If you're still experiencing issues:

1. Check the [GitHub Issues](https://github.com/your-org/calendar-api/issues)
2. Review the [Documentation](docs/)
3. Contact the development team
4. Submit a new issue with:
   - Detailed error message
   - Steps to reproduce
   - Environment information
   - Relevant logs 