# Architecture Rules

## Service Layer Rule (Hybrid Approach)

### Default: Pattern (with Service layer)
Sundin ang pattern by default - Handler -> Service -> Repository.

### Exception: Simple CRUD (no Service)
Kung simple passthrough lang, direkta sa Repository.

### Decision Rule
IF service adds value (logic, orchestration, side effects)
    -> Use service
ELSE IF service is pure passthrough
    -> Skip service, use repo directly

## Gumamit ng SERVICE kung:
- BUSINESS LOGIC - validation, computation, rules
- ORCHESTRATION - multiple repos, external calls
- SIDE EFFECTS - audit log, notifications, events
- TRANSACTIONS - multiple DB operations
- REUSABLE - across HTTP, CLI, jobs, gRPC, GraphQL

## Hindi kailangan ng SERVICE kung:
- Simple passthrough - Handler -> Repo -> Response
- Single repo, single operation - List, Get by ID
- Walang business rules - pure CRUD
- Hindi reusable - one handler only

## Examples sa RBAC
| Operation | Service? | Reason |
|---|---|---|
| List roles | No | Simple passthrough |
| Get role | No | Simple get |
| Create role | Yes | Duplicate check + validation |
| Update role | Yes | Validation |
| Delete role | Yes | Dependency check |
| AssignPermission | Yes | Orchestration (role + perm) |
| RevokePermission | No | Simple delete |
| ChangeUserRole | Yes | Audit + validation |

## Bakit Hybrid (hindi pure Service)

### Pros of Hybrid:
- Iwasan redundant - hindi lahat kailangan service
- Pragmatic - simple CRUD mabilis
- Complex logic - proper layering

### Pros of Pure Service:
- Consistency - lahat same structure
- Testability - proper unit tests
- Reusability - future-proof

### Recommendation:
Hybrid - balance ng dalawa. Pragmatic sa simple, proper sa complex.

## Industry Standard
Google, Uber, Netflix - ginagamit ang hybrid approach:
- Simple CRUD -> Handler -> Repo
- Complex logic -> Handler -> Service -> Repo

## Rule of Thumb
Every layer should add value.

If layer = passthrough -> remove it
If layer = adds value -> keep it

### Layer adds value if:
- May validation
- May transformation
- May orchestration
- May side effects
- May business rules

### Layer is redundant if:
- Pure passthrough
- Same signature as repo
- Walang added logic

## Flow Diagram
Handler (HTTP)
    |
    +-- Complex logic --> Service --> Repository
    |                                  |
    +-- Simple CRUD ---> Repository ---+
                                       |
                                    Database

## Examples

### Example 1: Simple CRUD (no service)
func (h *UserHandler) List(c echo.Context) error {
    users, err := h.userRepo.List(ctx)
    return c.JSON(200, users)
}

### Example 2: Complex logic (with service)
func (h *UserHandler) Register(c echo.Context) error {
    user, err := h.userService.Register(ctx, req)
    return c.JSON(201, user)
}

func (s *UserService) Register(ctx, req) (*User, error) {
    // 1. Validate email
    // 2. Check duplicate
    // 3. Hash password
    // 4. Create user
    // 5. Assign default role
    // 6. Send welcome email
    // 7. Audit log
    return user, nil
}

## Summary
| Rule | Detail |
|---|---|
| Default | Pattern - with Service layer |
| Exception | Simple CRUD - direct to Repo |
| Principle | Every layer must add value |
| Decision | Service if complex, Repo if simple |
| Consistency | Same pattern sa similar operations |
| Pragmatism | Don't add unnecessary layers |

Remember: Layering is a tool, not a religion. Use it where it adds value.
