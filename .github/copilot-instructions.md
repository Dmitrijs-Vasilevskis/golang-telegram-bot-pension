# Copilot PR review instructions

This repository is a Go Telegram bot. Focus reviews on correctness, simplicity, and consistency with existing code.

## What to check first (priority order)

1. **Correctness**
   - Does the code do what it claims?
   - Are edge cases handled (empty input, nil fields, duplicates, errors)?

2. **Safety**
   - No panics in normal flow
   - Proper error handling
   - No nil pointer risks (common with Telegram updates)

3. **Database logic**
   - Queries have deterministic ordering when needed
   - No unnecessary full scans or repeated queries in loops
   - For chat messages:
     - ordering should use `message_id` (not `created_at`) unless explicitly justified
     - trimming logic must not delete newer messages accidentally

4. **Architecture**
   - No business logic inside handlers if avoidable
   - No DB logic leaking into unrelated layers
   - No unnecessary new abstractions or layers

5. **Go best practices**
   - Idiomatic Go (simple, readable)
   - No unnecessary interfaces or abstractions
   - Proper use of `context.Context`
   - No goroutine leaks or unsafe concurrency

---

## Common mistakes to flag

- Missing error handling
- Ignoring returned errors
- Overly complex solutions where a simple one exists
- Introducing new patterns inconsistent with the repo
- Large refactors in small PRs

---

## What to avoid suggesting

- New frameworks or dependencies
- Large architectural changes
- Unnecessary abstractions
- Rewriting working code without clear benefit

---

## Expected suggestions

- Prefer small, targeted fixes
- Prefer clearer and simpler code
- Keep consistency with existing repository structure
- When uncertain, ask for clarification instead of assuming
