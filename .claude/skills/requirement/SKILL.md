name: requirement
description: Use when starting a new feature development - orchestrates the complete requirement-to-implementation workflow with mandatory review checkpoints. Triggers on new feature requests, product development, or any task requiring PRD/API/UI design before coding.
---

# ShopJoy 需求开发流程规范

## When to Use This Skill

- User requests a new feature or product development
- Starting any task that requires design before implementation
- User says "I want to add X feature" or "We need to build Y"
- Planning phase before coding begins

## Flow Overview

```
Phase 1          Phase 2              Phase 3          Phase 4
需求分析    →    设计阶段       →    计划阶段    →    实施阶段
```

## Phase 1: Requirement Analysis (需求分析)

### Execution

| Item | Content |
|-----|---------|
| **Agents** | `voltagent-biz:product-manager` + `shopify-expert` |
| **Output** | PRD document → `docs/domains/{domain}/{date}-{domain}-prd.md` |
| **Review** | Use `superpowers:brainstorming` to clarify requirements |

### Checklist

- [ ] Business scenarios complete
- [ ] User roles defined
- [ ] Edge cases covered
- [ ] Non-functional requirements (performance, security, i18n)

### Exit Criteria

- PRD document completed
- At least 2 rounds of review
- User explicitly approved

---

## Phase 2: Design Stage (设计阶段) - Parallel

### 2.1 API Design (Backend)

| Item | Content |
|-----|---------|
| **Agents** | `voltagent-core-dev:backend-developer` + `api-designer` |
| **Output** | Tech design → `docs/domains/{domain}/{date}-{domain}-design.md` |
| **Reviewer** | `backend-development:backend-architect` |

### 2.2 UI Design (Frontend)

| Item | Content |
|-----|---------|
| **Agent** | `voltagent-core-dev:ui-designer` |
| **Output** | UI design → `docs/domains/{domain}/design/{date}-{domain}-ui-design.md` |
| **Reviewers** | `voltagent-biz:product-manager` + `shopify-expert` |

### 2.3 Frontend Tech Design

| Item | Content |
|-----|---------|
| **Agent** | `voltagent-core-dev:frontend-developer` |
| **Prerequisite** | API design + UI design completed |
| **Output** | Frontend tech doc → `docs/domains/{domain}/design/{date}-{domain}-frontend-design.md` |

### Exit Criteria

- All design documents completed
- At least 2 rounds of review each
- User explicitly approved

---

## Phase 3: Integration Review & Planning

### 3.1 Document Integration Review

| Item | Content |
|-----|---------|
| **Agent** | `voltagent-biz:product-manager` |
| **Review Targets** | PRD + API design + UI design + Frontend tech design |
| **Focus** | Completeness, consistency, missing requirements |

### 3.2 Development Plan

| Item | Content |
|-----|---------|
| **Agent** | `superpowers:writing-plans` |
| **Output** | Dev plan → `docs/plans/{date}-{feature}-plan.md` |

### Exit Criteria

- All documents reviewed together
- No missing requirements
- Development plan approved
- User explicitly approved

---

## Phase 4: Implementation Stage

### 4.1 Code Development

| Item | Content |
|-----|---------|
| **Agent** | `superpowers:subagent-driven-development` or parallel agents |
| **Basis** | Development plan document |
| **Strategy** | Frontend/backend separation, parallel for independent modules |

### 4.2 Code Review

| Item | Content |
|-----|---------|
| **Agent** | `code-refactoring:code-reviewer` |
| **Focus** | Code quality, requirement compliance, standards adherence |

---

## Critical Principles

### MUST

| # | Principle |
|---|-----------|
| 1 | **Full-stack coverage**: Every requirement must consider Frontend + Backend + UI + Database |
| 2 | **Documentation first**: All design docs must be completed and approved before development |
| 3 | **Consistency**: Final implementation must match PRD and UI design exactly - no missing features |
| 4 | **Two-round review**: Each phase output requires at least 2 rounds of review |
| 5 | **User confirmation**: User must explicitly approve each phase before proceeding |

### MUST NOT

| # | Prohibition |
|---|-------------|
| 1 | Start coding before documents are approved |
| 2 | Only consider backend, ignoring frontend/UI |
| 3 | Skip review stages |
| 4 | Proceed to next phase without user confirmation |

---

## Quick Reference

```
Requirement → PM+Shopify Expert → PRD → Brainstorming → Review(2x) → User OK
                                        ↓
                  ┌─────────────────────┴─────────────────────┐
                  ↓                                           ↓
            Backend Dev+API Designer                   UI Designer
                  ↓                                           ↓
            Tech Design Doc                              UI Design Doc
                  ↓                                           ↓
            Backend Architect Review                  PM+Shopify Review
                  ↓                                           ↓
                  └─────────────────────┬─────────────────────┘
                                        ↓
                                  Frontend Dev
                                        ↓
                                Frontend Tech Doc
                                        ↓
                                  PM Integration Review
                                        ↓
                                  Development Plan
                                        ↓
                                  Code Development (parallel)
                                        ↓
                                  Code Review + Fix
```

---

## Document Naming Convention

```
{YYYY-MM-DD}-{domain}-{type}.md
```

| Type | Format | Example |
|------|--------|---------|
| PRD | `{date}-{domain}-prd.md` | `2026-03-24-order-prd.md` |
| Schema | `{date}-{domain}-schema.md` | `2026-03-24-user-schema.md` |
| UI Design | `{date}-{domain}-ui-design.md` | `2026-03-24-payment-ui-design.md` |
| Tech Design | `{date}-{name}-design.md` | `2026-03-22-sku-design.md` |
| Frontend Design | `{date}-{domain}-frontend-design.md` | `2026-03-24-order-frontend-design.md` |

---

## Phase Checklist Template

Use this checklist at the end of each phase:

```
□ Output delivered to correct directory
□ Output naming follows convention
□ At least 2 rounds of review completed
□ User explicitly approved
□ docs/README.md index updated (if needed)
```