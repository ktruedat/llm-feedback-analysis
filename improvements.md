# Future Improvements & Enhancements

---

## 🔴 High Priority

### 1. **Add Unit and Integration Tests**

**Status**: Not implemented (time constraints)

**Why it matters**:

- Prevents regressions when adding features
- Documents expected behavior
- Enables confident refactoring
- Required for production systems

**Suggested approach**:

**Unit tests**:

```go
// Test domain logic (no dependencies)
func TestFeedbackBuilder_ValidRating(t *testing.T) {
    fb := feedback.NewBuilder().
        WithRatingValue(5).
        WithCommentText("Great!").
        Build()
    // Assert expectations
}
```

**Integration tests**:

```go
// Test repository layer with test database
func TestFeedbackRepository_Create(t *testing.T) {
    // Use testcontainers for PostgreSQL
    // Test actual database operations
}
```

**API tests**:

```go
// Test HTTP handlers end-to-end
func TestCreateFeedbackHandler(t *testing.T) {
    // Setup test server
    // Make HTTP request
    // Assert response
}
```

**Coverage goals**:

- Domain layer: 90%+ (pure business logic)
- Services: 80%+ (orchestration)
- Handlers: 70%+ (HTTP layer)
- Repositories: 60%+ (database operations)

**Tools**:

- `testing` package (Go standard library)
- `testify` for assertions
- `testcontainers-go` for PostgreSQL in tests
- `httptest` for HTTP testing

---

### 2. **Production-Ready Admin Account Creation**

**Current state**: Default admin created via migration with hardcoded credentials

**Security concerns**:

- ❌ Credentials in migration file (committed to git)
- ❌ Same password across all deployments
- ❌ No way to rotate credentials without database modification
- ❌ Not suitable for production

**Proposed solution**:

#### Option A: First-Time Setup Wizard

```go
// On first startup, check if admin exists
func (app *App) initializeDefaultAdmin(ctx context.Context) error {
    adminExists, err := app.userRepo.AdminExists(ctx)
    if err != nil {
        return err
    }
    
    if !adminExists {
        // Generate secure random password
        tempPassword := generateSecurePassword(32)
        
        // Create admin with temporary password
        admin, err := app.userService.CreateAdmin(ctx, "admin@temp.local", tempPassword)
        if err != nil {
            return err
        }
        
        // Output credentials to stdout (captured in logs)
        app.logger.Warn("=== TEMPORARY ADMIN CREDENTIALS ===")
        app.logger.Warn("Email: admin@temp.local")
        app.logger.Warn("Password: " + tempPassword)
        app.logger.Warn("CHANGE THESE CREDENTIALS IMMEDIATELY")
        app.logger.Warn("===================================")
        
        // Optionally: Save to Kubernetes secret, AWS Secrets Manager, etc.
        if secretsManager != nil {
            secretsManager.Store("admin-credentials", map[string]string{
                "email":    "admin@temp.local",
                "password": tempPassword,
            })
        }
    }
    
    return nil
}
```

#### Option B: Environment Variable Setup

```bash
# On first deployment, set these environment variables
INITIAL_ADMIN_EMAIL=admin@company.com
INITIAL_ADMIN_PASSWORD=$(openssl rand -base64 32)
FORCE_PASSWORD_CHANGE=true
```

#### Option C: CLI Command

```bash
# Backend provides admin creation command
./backend admin create --email admin@company.com --role admin

# Outputs:
# Admin created successfully
# Email: admin@company.com  
# Temporary password: [generated-password]
# Password expires in 24 hours
```

**Best practice**:

- Use secrets management (AWS Secrets Manager, HashiCorp Vault, Kubernetes Secrets)
- Force password change on first login
- Set password expiration
- Remove migration-based admin creation

---

### 3. **Cache OpenAI System Prompt**

**Current state**: System prompt sent with every analysis request

**Opportunity**: OpenAI supports prompt caching to reduce costs

**Implementation**:

```go
// openai.go
func (c *OpenAIClient) buildRequestBody(userPayload Map) ([]byte, error) {
systemPrompt := c.buildSystemPrompt()
    systemPrompt := c.buildSystemPrompt()
    
    requestBody := Map{
        "model": c.model,
        "input": []Map{
            {
                "role":    "system",
                "content": systemPrompt,
                "cache_control": Map{  // Enable caching
                    "type": "ephemeral",
                },
            },
            {
                "role":    "user",
                "content": string(userJSON),
            },
        },
        // ... rest of request
    }
    
    }
```

**Cost savings**:

- System prompt: ~500 tokens
- Cached input cost: $0.025 per 1M tokens (vs $0.25 regular)
- **90% reduction** on system prompt tokens
- Especially valuable with frequent analyses

**Calculation**:

```
Without caching:
- 1,000 analyses × 500 tokens = 500,000 tokens
- Cost: 500,000 × $0.25/1M = $0.125

With caching:
- First request: 500 tokens × $0.25/1M = $0.000125
- Next 999 requests: 499,500 tokens × $0.025/1M = $0.012
- Total: $0.012 (90% savings)
```

**References**:

- [OpenAI Prompt Caching](https://platform.openai.com/docs/guides/prompt-caching)

---

### 4. **Orphaned Feedback Retry Logic**

**Current state**: System tracks orphaned feedbacks but doesn't retry them

**Problem**: If analysis fails after OpenAI call but before database save:

- Feedbacks marked as analyzed
- No analysis results stored
- Lost insights, wasted API costs

**Existing infrastructure** (already in place):

```sql
-- feedback_analysis_assignments table tracks which feedbacks were analyzed
CREATE TABLE feedback.feedback_analysis_assignments
(
    feedback_id UUID,
    analysis_id UUID,
    created_at  TIMESTAMP
);
```

**Proposed solution**:

#### A. Startup Recovery Check

```go
// analyzer.go
func (a *analyzer) recoverOrphanedFeedbacks(ctx context.Context) error {
logger := a.logger.WithSpan(ctx)
    logger := a.logger.WithSpan(ctx)
    
    // Find feedbacks that should have been analyzed but aren't in any assignment
    orphaned, err := a.analysisRepo.FindOrphanedFeedbacks(ctx)
    if err != nil {
        return fmt.Errorf("failed to find orphaned feedbacks: %w", err)
    }
    
    if len(orphaned) == 0 {
        logger.Info("no orphaned feedbacks found")
        return nil
    }
    
    logger.Warn("found orphaned feedbacks, triggering recovery analysis", 
        "count", len(orphaned))
    
    // Re-analyze orphaned feedbacks
    for _, fb := range orphaned {
        a.EnqueueFeedback(fb)
    }
    
    }
```

#### B. Periodic Health Check

```go
// Run health check every hour
func (a *analyzer) startHealthCheck(ctx context.Context) {
ticker := time.NewTicker(1 * time.Hour)
    ticker := time.NewTicker(1 * time.Hour)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            if err := a.recoverOrphanedFeedbacks(ctx); err != nil {
                a.logger.Error("orphaned feedback recovery failed", err)
            }
        case <-ctx.Done():
            return
        }
    }
```

#### C. Repository Method

```go
// analysis/repository.go
func (r *repo) FindOrphanedFeedbacks(ctx context.Context) ([]*feedback.Feedback, error) {
    query := `
        SELECT f.* FROM feedback.feedbacks f
        WHERE f.analyzed = true
        AND NOT EXISTS (
            SELECT 1 FROM feedback.feedback_analysis_assignments faa
            WHERE faa.feedback_id = f.id
        )
        ORDER BY f.created_at ASC
        LIMIT 100;  -- Process in batches
    `
    // Execute query and map to domain entities
}
```

**Benefits**:

- Automatic recovery from transient failures
- No data loss from network partitions
- Maintains analysis completeness
- Uses existing database schema

---

## 🟡 Medium Priority

### 5. **Health Check & Readiness Endpoints**

**Current state**: No health checks

**Proposed endpoints**:

```go
// GET /health - Basic liveness probe
func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "status":    "healthy",
        "timestamp": time.Now().Format(time.RFC3339),
    })
}

// GET /ready - Readiness probe (checks dependencies)
func (h *Handlers) Ready(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
    // Check database
    if err := h.db.Ping(ctx); err != nil {
        w.WriteHeader(http.StatusServiceUnavailable)
        json.NewEncoder(w).Encode(map[string]string{
            "status": "not ready",
            "reason": "database unavailable",
        })
        return
    }
    
    // Check OpenAI connectivity (optional)
    if err := h.llmClient.HealthCheck(ctx); err != nil {
        w.WriteHeader(http.StatusServiceUnavailable)
        json.NewEncoder(w).Encode(map[string]string{
            "status": "degraded",
            "reason": "LLM service unavailable",
        })
        return
    }
    
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "status": "ready",
    })
}
```

**Kubernetes integration**:

```yaml
livenessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 10
  periodSeconds: 30

readinessProbe:
  httpGet:
    path: /ready
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 10
```

---

### 6. **Metrics (Prometheus)**

**Current state**: Only distributed tracing (no metrics)

**Proposed additions**:

```go
// metrics/metrics.go
var (
    FeedbacksCreated = promauto.NewCounter(prometheus.CounterOpts{
        Name: "feedbacks_created_total",
        Help: "Total number of feedbacks created",
    })
    
    AnalysesTriggered = promauto.NewCounter(prometheus.CounterOpts{
        Name: "analyses_triggered_total",
        Help: "Total number of analyses triggered",
    })
    
    AnalysisDuration = promauto.NewHistogram(prometheus.HistogramOpts{
        Name:    "analysis_duration_seconds",
        Help:    "Time taken to complete analysis",
        Buckets: prometheus.DefBuckets,
    })
    
    OpenAICost = promauto.NewCounter(prometheus.CounterOpts{
        Name: "openai_tokens_used_total",
        Help: "Total OpenAI tokens consumed",
    })
)
```

**Usage**:

```go
// In analyzer service
start := time.Now()
result, err := c.llmClient.AnalyzeFeedbacks(ctx, feedbacks, previousAnalysis)
AnalysisDuration.Observe(time.Since(start).Seconds())
OpenAICost.Add(float64(result.TokensUsed))
AnalysesTriggered.Inc()
```

**Grafana dashboard**:

- Feedbacks per minute
- Analysis latency (p50, p95, p99)
- OpenAI token usage trends
- Cost projection based on usage

---

### 7. **Background Job Queue (e.g., NATS/Redis)**

**Current state**: Analysis uses simple goroutine + channel

**Limitation**:

- No persistence if server restarts mid-analysis
- No visibility into job status
- No retry configuration
- Can't scale across multiple instances

**Proposed solution** (using Asynq):

```go
// jobs/analysis_job.go
type AnalysisJob struct {
    FeedbackIDs []string `json:"feedback_ids"`
}

func (j *AnalysisJob) ProcessTask(ctx context.Context, task *asynq.Task) error {
    var payload AnalysisJob
    if err := json.Unmarshal(task.Payload(), &payload); err != nil {
        return err
    }
    
    // Fetch feedbacks and run analysis
    feedbacks, err := fetchFeedbacks(ctx, payload.FeedbackIDs)
    if err != nil {
        return err
    }
    
    return runAnalysis(ctx, feedbacks)
}

// Enqueue job
func EnqueueAnalysis(client *asynq.Client, feedbackIDs []string) error {
    payload, _ := json.Marshal(AnalysisJob{FeedbackIDs: feedbackIDs})
    task := asynq.NewTask("analysis:run", payload)
    
    return client.Enqueue(task, 
        asynq.MaxRetry(3),
        asynq.Timeout(5*time.Minute),
    )
}
```

**Benefits**:

- Job persistence (survives restarts)
- Built-in retry logic
- Distributed processing across multiple workers
- Dashboard for monitoring jobs
- Priority queues

---

## 🟢 Nice to Have

### 8. **Redis Caching for Analysis Results**

**Use case**: Admin dashboard shows same analysis multiple times

```go
// Cache latest analysis for 5 minutes
func (s *service) GetLatestAnalysis(ctx context.Context) (*analysis.Analysis, error) {
cacheKey := "analysis:latest"
    cacheKey := "analysis:latest"
    
    // Try cache first
    if cached, err := s.cache.Get(ctx, cacheKey); err == nil {
        var result analysis.Analysis
        json.Unmarshal([]byte(cached), &result)
        return &result, nil
    }
    
    // Cache miss - fetch from database
    result, err := s.repo.GetLatest(ctx)
    if err != nil {
        return nil, err
    }
    
    // Cache for 5 minutes
    data, _ := json.Marshal(result)
    s.cache.Set(ctx, cacheKey, data, 5*time.Minute)
    
    }
```

---

### 9. **Feedback Search & Filtering**

**Current**: List all feedbacks with pagination only

**Proposed**:

```http
GET /api/v1/feedbacks?search=performance&rating=1,2&from=2024-01-01&to=2024-12-31
```

---

### 10. **Webhook Notifications**

**Trigger**: When analysis completes

```http
POST https://customer-webhook.com/analysis-complete
Content-Type: application/json

{
  "analysis_id": "uuid",
  "overall_sentiment": "positive",
  "feedback_count": 42
}
```

---

### 11. **Export Functionality**

**Formats**: CSV, PDF, JSON

```http
GET /api/v1/analyses/{id}/export?format=
GET /api/v1/analyses/{id}/export?format = pdf
```
