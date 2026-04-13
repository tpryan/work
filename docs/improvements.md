This review focuses on improving the **idiomatic structure**, **maintainability**, and **documentation** of the `work` tool (EtchWake). The codebase has a solid functional foundation but exhibits several patterns—such as inconsistent package responsibilities and non-standard error handling—that can be refined for a more "Go-like" experience.

### 1. Refine Package Responsibilities & Boundaries
The current structure has some "leaky abstractions" where internal details are exposed or duplicated across packages.

* **Move Business Logic out of `internal`:** In Go, `internal` is for code you don't want others to import. Since `Config` and `Artifact` are the core of your application, consider moving the logic in `internal/artifact` directly into the top-level `work` package. This simplifies imports: instead of `artifact.Artifact`, you use `work.Artifact`.
* **Consolidate Auth Logic:** You have authentication logic split between `internal/googleauth` and `internal/option`. Consolidate these into a single `pkg/auth` or keep them in `internal/auth`.
* **Interface-Driven Providers:** Instead of the `main.go` calling `github.Search` and `drive.Search` directly, define a `Source` interface:
    ```go
    type Source interface {
        Fetch(ctx context.Context) (Artifacts, error)
    }
    ```
    This allows you to add new sources (like a new internal tool) without touching the core reporting logic.

### 2. Idiomatic Go Improvements
* **Standardize Date Handling:** In `work.go`, `Criteria` uses `time.Time`, but the YAML unmarshaler may struggle with various formats. Implement `UnmarshalYAML` for a custom `Date` type to handle the `YYYY-MM-DD` format more robustly.
* **Context Propagation:** Many functions in `internal/gsheet` and `internal/github` do not accept a `context.Context`. To make the tool "production-ready," pass `ctx` through all network-calling functions to allow for timeouts and cancellations.
<!-- * **Error Wrapping:** Currently, some errors are returned directly or formatted with `fmt.Errorf("...: %s", err)`. Use the `%w` verb (e.g., `fmt.Errorf("failed to fetch: %w", err)`) to allow callers to use `errors.Is` or `errors.As`. -->
<!-- * **Avoid `log.Fatalf` in Libraries:** Functions in `internal/` should return errors to the caller rather than terminating the program with `log.Fatalf`. Reserve `Fatalf` for the `main` package. DONE -->

### 3. Documentation & API Design
* **Replace `Massage` with Functional Options:** The `Massage` function name is non-standard. A more idiomatic approach for filtering/transforming slices in Go is to use a "Filter" or "Pipeline" pattern.
* **Use Go Doc Comments:** While some functions have comments, many are missing or don't follow the `// Name ...` convention.
    * *Example:* `// Artifacts represents a collection of work products.` instead of just `// Artifacts is a collection...`.
* **README and Examples:** The `docs/configuration.md` is excellent. Consider adding a `README.md` at the root that explains how to build the project using the existing `Makefile`.

### 4. Enhancing Performance & Safety
<!-- * **Concurrency Control:** In `writeReport`, you launch a goroutine for every destination. While efficient for a few sheets, if the list grows, you might hit API rate limits. Consider using a worker pool or a `semaphore` to limit concurrent API writes.
* **Struct Tags:** Ensure all fields in `Artifact` that are intended for Sheet output are clearly mapped. You currently use a slice of interfaces for Sheet updates, which is manual and error-prone. Consider a reflection-based approach or a dedicated `Marshaler` to convert `Artifact` to row data. -->

### 5. Code Cleanliness (The "Mess")
* **Dead Code:** `internal/googleauth/client.go` contains a lot of commented-out reasoning and redundant `NewClientOption` logic. Clean these up to make the intent clear.
* **Consistency:** The `.gitignore` suggests many user-specific YAML files are being tracked or ignored individually. Use a pattern like `users/*.yaml` to keep the root clean.
* **Makefile:** The current `Makefile` is very minimal. Add targets for `lint`, `fmt`, and `build` to automate the cleanup of the "mess" during development.