# Architecture Overview

## System Architecture

```mermaid
graph TB
    subgraph Client
        C[API Client]
        CLI[CLI Tools]
        SDK[SDK/Client Libraries]
    end

    subgraph AWS Cloud
        subgraph VPC
            subgraph Public Subnet
                ALB[Application Load Balancer]
                NAT[NAT Gateway]
            end

            subgraph Private Subnet
                subgraph ECS Cluster
                    API1[Calendar API Container 1]
                    API2[Calendar API Container 2]
                    API3[Calendar API Container 3]
                end

                subgraph RDS
                    DB[(PostgreSQL Primary)]
                    REP[(Read Replica)]
                end
            end
        end

        subgraph Supporting Services
            ECR[ECR Repository]
            S3[S3 Bucket]
            Secrets[Secrets Manager]
            CloudWatch[CloudWatch]
            ACM[Certificate Manager]
        end
    end

    C -->|HTTPS| ALB
    CLI -->|HTTPS| ALB
    SDK -->|HTTPS| ALB
    ALB -->|HTTP| API1
    ALB -->|HTTP| API2
    ALB -->|HTTP| API3
    API1 -->|TCP/5432| DB
    API2 -->|TCP/5432| DB
    API3 -->|TCP/5432| DB
    API1 -->|TCP/5432| REP
    API2 -->|TCP/5432| REP
    API3 -->|TCP/5432| REP
    API1 -->|Read/Write| Secrets
    API2 -->|Read/Write| Secrets
    API3 -->|Read/Write| Secrets
    ECR -->|Pull| API1
    ECR -->|Pull| API2
    ECR -->|Pull| API3
    S3 -->|State| Terraform
    API1 -->|Logs/Metrics| CloudWatch
    API2 -->|Logs/Metrics| CloudWatch
    API3 -->|Logs/Metrics| CloudWatch
    ALB -->|SSL Cert| ACM
```

## Component Description

### Client Layer
- **API Client**: Any HTTP client that can make REST API calls
- **CLI Tools**: Command-line interface for automation
- **SDK/Client Libraries**: Language-specific client libraries
- **Authentication**: Uses API key via X-API-Key header

### Load Balancer
- **Application Load Balancer (ALB)**
  - SSL/TLS termination via ACM
  - Health checks with configurable thresholds
  - Request routing with path-based rules
  - Security groups with least privilege
  - WAF integration for DDoS protection

### Application Layer
- **ECS Container**
  - Go application with optimized runtime
  - Environment configuration via Secrets Manager
  - Health monitoring with CloudWatch
  - Structured logging with correlation IDs
  - Auto-scaling based on metrics
  - Container insights for performance

### Database Layer
- **RDS PostgreSQL**
  - Managed database service with automated maintenance
  - Automated backups with point-in-time recovery
  - High availability with multi-AZ deployment
  - Read replicas for scaling
  - Security groups with encrypted connections
  - Performance insights and monitoring

### Supporting Services
- **ECR**: Container registry with image scanning
- **S3**: Terraform state storage with versioning
- **Secrets Manager**: Secure credential storage with rotation
- **CloudWatch**: Centralized logging and monitoring
- **ACM**: SSL certificate management

## Security Architecture

```mermaid
graph TB
    subgraph Security
        subgraph Network Security
            SG[Security Groups]
            NACL[Network ACLs]
            VPC[VPC Isolation]
            WAF[Web Application Firewall]
            SHIELD[Shield Advanced]
        end

        subgraph Application Security
            SSL[SSL/TLS]
            API[API Key Auth]
            CORS[CORS Policy]
            RATE[Rate Limiting]
            VALID[Input Validation]
        end

        subgraph Data Security
            ENC[Encryption at Rest]
            TRANS[Encryption in Transit]
            BACK[Automated Backups]
            AUDIT[Audit Logging]
            ROTATE[Key Rotation]
        end

        subgraph Compliance
            ISO[ISO 27001]
            SOC[SOC 2]
            GDPR[GDPR]
            HIPAA[HIPAA]
        end
    end
```

## Deployment Architecture

```mermaid
graph TB
    subgraph CI/CD Pipeline
        subgraph GitHub Actions
            subgraph Quality Gates
                LINT[Code Linting]
                TEST[Unit Tests]
                SEC[Security Scan]
                COV[Coverage Check]
            end

            subgraph Build Phase
                BUILD[Build Container]
                SCAN[Container Scan]
                PUSH[Push to ECR]
            end

            subgraph Deploy Phase
                TF_PLAN[Terraform Plan]
                TF_APPLY[Terraform Apply]
                DB_MIG[Database Migrations]
                HEALTH[Health Check]
            end
        end

        subgraph Environments
            DEV[Development]
            STG[Staging]
            PRD[Production]
        end

        subgraph Monitoring
            ALERT[Alerting]
            DASH[Dashboards]
            LOGS[Log Aggregation]
        end
    end

    LINT --> TEST
    TEST --> SEC
    SEC --> COV
    COV --> BUILD
    BUILD --> SCAN
    SCAN --> PUSH
    PUSH --> TF_PLAN
    TF_PLAN --> TF_APPLY
    TF_APPLY --> DB_MIG
    DB_MIG --> HEALTH
    HEALTH --> DEV
    DEV --> STG
    STG --> PRD
    PRD --> ALERT
    ALERT --> DASH
    DASH --> LOGS
```

## Monitoring Architecture

```mermaid
graph TB
    subgraph Monitoring
        subgraph Application Metrics
            HC[Health Checks]
            LOG[Application Logs]
            MET[Custom Metrics]
            TRACE[Distributed Tracing]
            PROF[Profiling]
        end

        subgraph Infrastructure Metrics
            CPU[CPU Usage]
            MEM[Memory Usage]
            NET[Network Traffic]
            DISK[Disk Usage]
            CONN[Connection Count]
        end

        subgraph Database Metrics
            CONN[Connections]
            QUERY[Query Performance]
            SIZE[Database Size]
            REPL[Replication Lag]
            CACHE[Cache Hit Ratio]
        end

        subgraph Business Metrics
            REQ[Request Rate]
            ERR[Error Rate]
            LAT[Latency]
            USAGE[API Usage]
            COST[Cost Metrics]
        end
    end

    subgraph Alerting
        ALERT[Alert Manager]
        PAGER[PagerDuty]
        SLACK[Slack]
        EMAIL[Email]
    end

    HC --> ALERT
    LOG --> ALERT
    MET --> ALERT
    TRACE --> ALERT
    PROF --> ALERT
    CPU --> ALERT
    MEM --> ALERT
    NET --> ALERT
    DISK --> ALERT
    CONN --> ALERT
    QUERY --> ALERT
    SIZE --> ALERT
    REPL --> ALERT
    CACHE --> ALERT
    REQ --> ALERT
    ERR --> ALERT
    LAT --> ALERT
    USAGE --> ALERT
    COST --> ALERT
    ALERT --> PAGER
    ALERT --> SLACK
    ALERT --> EMAIL
```

## Data Flow

1. Client sends authenticated request to ALB
2. ALB performs SSL termination and WAF checks
3. Request is routed to healthy ECS container
4. Application validates API key and request
5. Database operations performed with connection pooling
6. Response is logged and metrics are recorded
7. Response returned to client with appropriate headers

## High Availability

- Multi-AZ RDS deployment with automated failover
- ECS container redundancy across AZs
- ALB health checks with configurable thresholds
- Automated failover with zero downtime
- Read replicas for database scaling
- Connection pooling for efficient resource usage

## Disaster Recovery

- Automated database backups with point-in-time recovery
- Terraform state versioning in S3
- Container image versioning in ECR
- Infrastructure as Code with modular design
- Cross-region replication for critical data
- Automated recovery procedures

## Scaling

- Horizontal scaling via ECS with auto-scaling groups
- Database read replicas for read scaling
- Connection pooling for efficient resource usage
- Caching strategies with Redis (optional)
- Load balancing with health checks
- Resource-based auto-scaling policies

## Cost Optimization

- Spot instances for non-critical workloads
- Reserved instances for predictable loads
- Auto-scaling based on demand
- Resource tagging for cost allocation
- CloudWatch cost monitoring
- Regular cost reviews and optimization 