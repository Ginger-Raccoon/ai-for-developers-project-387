# Contributing

## Commit format — Conventional Commits

All commits (including commits made by AI agents) **must** follow [Conventional Commits](https://www.conventionalcommits.org/).

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

### Types

| Type       | When to use                                      | Version bump |
|------------|--------------------------------------------------|--------------|
| `feat`     | New feature visible to users or consumers of API | minor        |
| `fix`      | Bug fix                                          | patch        |
| `feat!`    | Breaking change (or `BREAKING CHANGE:` in footer)| major        |
| `fix!`     | Breaking bug fix                                 | major        |
| `perf`     | Performance improvement                          | patch        |
| `refactor` | Code change that is not a fix or feature         | —            |
| `test`     | Adding or fixing tests                           | —            |
| `ci`       | CI/CD configuration changes                      | —            |
| `docs`     | Documentation only                               | —            |
| `chore`    | Maintenance (deps update, tooling, etc.)         | —            |
| `style`    | Formatting, missing semicolons, etc.             | —            |

### Scopes (optional but recommended)

| Scope      | Area                     |
|------------|--------------------------|
| `backend`  | Go backend               |
| `frontend` | React frontend           |
| `api`      | TypeSpec / OpenAPI spec  |
| `ci`       | GitHub Actions workflows |

### Examples

```
feat(backend): add slot generation for custom durations
fix(frontend): show error message when booking slot is taken
feat(api): add 404 response to POST /bookings
test(backend): add integration tests for double-booking scenario
ci: configure release-please workflow
chore(frontend): update antd to v5.24
feat!: rename eventTypeId to eventType in booking response
```

### Breaking changes

Add `!` after the type or include `BREAKING CHANGE:` in the footer:

```
feat!(api)!: rename field guestName to name in CreateBookingRequest

BREAKING CHANGE: the field `guestName` in CreateBookingRequest is now `name`.
Frontend and backend must be updated together.
```

## Release process

Releases are automated via **release-please**.

1. Merge commits with `feat` or `fix` types into `main`.
2. release-please automatically opens or updates a **release PR** with:
   - Bumped version in `version.txt`
   - Generated `CHANGELOG.md` entries
3. Merge the release PR → a GitHub Release is created automatically.

`feat` commits → minor version bump (`0.1.0` → `0.2.0`)  
`fix` commits → patch version bump (`0.1.0` → `0.1.1`)  
Breaking changes → major version bump (`0.1.0` → `1.0.0`)
