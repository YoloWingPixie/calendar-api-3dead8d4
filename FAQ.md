# FAQ

## What would you have done differently?

1. I started with Python in this project, however I did strictly typed python with very pedantic pydantic configurations. I should have just started with golang. Golang is: 1) Faster, 2) statically typed by default. For some reason, I envisioned that a FastAPI/SQLAlchemy/Pydantic setup with Alembic for migrations would take no time at all, that tended to be very much false and I could have developed this much quicker just starting with Go.

2. I have this idea when I set out to make any API-shaped object project: Needs to have dynamic route loading, needs to have a solid config class with a strong data model, needs to be highly modular, authentication should be a first class member of the project features. Frankly, there is a minimum amount of time it takes to do all of this and AI sometimes slows you down in this regard because frankly it's terrible at systems level thinking. Had I just made a toy application from the start, this would've taken me less time.

## How did you use AI in this process?

1. Planning: [Gemini AI Studio \ Cursor Tab ]
I strongly believe that with AI-assisted coding there is nothing more important than building out a complete documentation package BEFORE putting a single line of code to source. For that reason, I start with the domain model. What are my entities, what are my state machines (if applicable), how do entities relate to each other, what are my invariants. For this project, I then went through and did all of my ADRs, and determined exactly what I was going to do. This, still changed, but was important for keeping AI on task. Then I made my data model, then I made a represantive OpenAPI spec for what I wanted the API to do. I made my assumptions, and then finally I created the PRD. 

I heavily use Deep Research to investigate ADR decisions and implications, and then, I use Gemini AI Studio to assist writing the ADR to a template spec, and then I red team the ADR with a separate session. I do not necessarily agree with everything on both ends, but I like to see the additional takes. Sometimes, there is something relevant. I had strong opinions about the ADRs already, before the project, so I was pretty much aware of the choices and the implementations before hand, but I did use this process to refine my thoughts and show the exercise.

2. Execution: [Claude Code \ Cursor Tab \ Gemini AI Studio]
Most of the project was done by:

1. Telling Claude Code to constantly keep docs in context and then pick off one work package from the PRD at a time. 
2. Using Cursor Tab to correct Claude Code's mistakes, and to add additional context to the code.
3. In some cases, I would do refactoring through AI Studio when I had a very specific and exact plan for what I wanted that I wanted to ensure was done with no context pollution and by a very low temperature model. 
4. In some cases, I would write the first example for Claude Code to show how I wanted things done idiomatically, and then I would ask it to continue on that work package. 
5. In almost all cases I would write the unit tests in a separate session and with no context of the source code. This generally reduces the amount of gaming an AI does to get the tests to pass or excessive mocking. I would then use Cursor Tab to correct the tests and add additional context.

## What would you add to this project or change it to productionalize it?
1. I would, had I had more time, used go-swagger for the actual generation of the server code. I didn't do that, because frankly, I forgot about it, wrote out the go server, and then when I realized without FastAPI I didn't have a nice endpoint to just serve the OpenAPI spec, I looked around and remembered go-swagger existed. "But did you try to have AI implement it?" Yes, I did, but it didn't work out in one-shot and I did not want to risk the time going down the rabbit hole. Golang was already a late decision after becoming fed up with the fragility of the Python codebase, so I just wanted to get it done.

2. There is a notional user authentication system with tokens, but you can't add users through the API. I would have added Cognito or some other IDP to the project to allow for user management and authentication.

3. The project *does* have an ACM, but it is self signed to the default URL of the ALBs, which is obviously not ideal and throws browser warnings. We would need to add proper custom domain management and ACM certificates to the ALBs for those domains in terraform. 

4. The project places all of the environments in the same AWS account, which is poor security hygeine. I would have split the environments into separate AWS accounts, and then used a cross-account role to allow the CI/CD pipeline to deploy to the other accounts. This would also require Doppler to sync secrets to tohe other accounts. Ephemeral environments would then go in the dev account.

5. I demonstrate that there is a PR environment cleanup, but there aren't actually many situations in which an Ephemeral environment would be created, since PRs to main use staging.

6. I would add a more formalized tool for VERSION management instead of the shell grepping that I am doing currently.

7. I would add notional Github releases to the project, for the sake of having a release history. I would have it autogenerate release notes based on the KACL format in CHANGELOG.md.

8. I would provide a generated client library for the API using go-swagger.

9. The metrics endpoint is unauthenticated, but nothing reads it in this project. For production we would obviously want to secure it, either through an API token or host based authentication.

10. On that matter, there is the batteries included for telemetry, but there is no actual telemetry service in use, aside from CloudWatch. That is observability, but I would want to have a more advanced, centralized telemetry service like Grafana in use.

11. There's no rate limiting, there should be.

12. On the subject of scaling to prod, based on the assumptions that I made, I don't think the current implementation has many issues scaling to prod based on the expected RPS. ECS already handles deployments and scaling, and the database in RDS is already set up to scale. ECS allows for surging so during deployments there isn't currently any downtime. Obviously since we only have 1 task, there's no advanced deployment strategy that is applicable. ECS does have the ability to setup blue green deployments with traffic control, so I would probably use that if we determined that more than 1-2 tasks was required for ECS. 

13. Does a calendar-api need to be multi-region if its internal? I don't think so, but we could always utilize a multi-region RDS setup with Aurora, and then perhaps use Route53's latency based routing to route users to the closest region, in a region down scenario this would *probably* work for failover but if we need something more formal for failover then we could use Global Accelerator for failover, but that is probably way too expensive for this.

## Ok, but you low balled yourself in the assumptions by saying this was an internal tool used by a company whose main business activity isn't sending calendar invites. What if this was a public API that had to scale to millions of users?
1. I would use Read Replicas for RDS. Obviously we would have the concept of tenants, calendars, and full users in this scenario, so we would need to hav ea multi-tenant database schema. I might consider having a discrete database for events / calendars with its own read replicas since those are "global" resources and then had users and tenants be its own separate auth database with read replicas (if the RPS was high enough). 

2. I would use a caching layer like Redis to cache the events and calendars, since those are the most read resources

3. I would probably just go ahead and use Kubernetes for the compute workload, and then use a global traffic load balancer to failover between clusters in different regions. ECS Fargate is great on the lower end of things, but when you have a large RPS, then the Kubernetes control plane cost becomes so negligible, an idiomatic k8s solution is often less work as long as teams are disciplined and you're using a good service mesh like istio. 

4. I would keep things boring, good scalable architecture is boring and has very few tricks. anything beyond the above would need to be an emergent requirement based on actual notable issues. 

## What Alerts and Monitors would you use?
1. We already have a health check for the container/API itself. 
2. error rate / requst rate % over X time frame (perhaps 15 minutes), my gut says anything above 4% is probably strange but that really depends on volume. We would probably want to restrict this to 500s and non 401s/404s.
3. We could still have a 404 monitor for unusually high amounts of 404s if desired
4. Latency 95p above 400ms would likely be concerning.
5. Availability monitoring for SLO/SLA validation
6. Database metrics like locks, read, writes, latency
7. If this were a global publically available API I would also implement a monitor and alert for traffic falling off more than 33% Day over Day during the same time period per day, that could often point to a regional networking issue or another availaiblity issue that is regionalas

## Security Improvements
1. Probably could use a better CORS implementation, if more was known about how this API should be accessed.
2. Would have users
3. API Keys could be forced to expire with scopes
4. I would have permissions on calendars, events would need to associate to a calendar, and users would need to be granted access to calendars other than their personal calendar.
5. As mentioned in assumptions, we could consider that title and description may contain PII or otherwise classified information, in a privacy focused context that could mean that only participants in meetings using the API could acutally view the title or description. It does mean we should consider data at rest encryption as well on those columns.
6. There's basically no input validation. We do enforce the start time must be before end time invariant, but we don't enforce things like username length, title lenght, and description. Those fields are also not *properly* sanitized.I did add some sanitization to show that I was thinking about that and database queries are parameterized to prevent *some* SQL injection but I would be lying if I said I thought that deeply about it. 