# Ambiguities, Assumptions, and Possible Improvements

## Ambiguities

1. Should `robots.txt` be respected when checking links?
2. How should "internal" vs. "external" links be defined?
3. How to detect a "login form"?
4. How should HTML version be detected (e.g., DOCTYPE vs. heuristics)?
5. What is the expected performance for checking large numbers of links?
6. How should invalid or malformed URLs be handled in user input?
7. What are the UI/UX requirements for displaying results?
8. Should redirects (3xx) be considered accessible or inaccessible?
9. How should headings be counted (visible vs. hidden elements)?

## Assumptions

1. For simplicity, I haven't implemented `robots.txt` parsing
2. Any link with the same `Hostname` is considered "internal"; others are "external"
3. Using password input and keywords in forms to detect login forms
4. Relying on DOCTYPE for HTML version detection
5. TBD
6. Input URLs must start with `http://` or `https://`; otherwise, default to `https://`
7. A simple HTML interface with Go templates
8. 3xx redirects are treated as accessible
9. All headings (h1-h6) in the DOM are counted, regardless of visibility

## Possible Improvements

- [] Store analysis results in a database for historical tracking
- [] AuthN and AuthZ?
- [] Use cache to improve performance for repeated analyses
- [] Use an SPA framework (e.g., React, Vue) for a more dynamic UI
- [] CD pipeline for automated deployments to cloud providers (e.g., Azure AKS, AWS EKS) using GitHub Actions and Terraform
- [] Add metrics, tracing, and logging (e.g., Grafana Stack and OpenTelemetry)
- [] Security and vulnerability enhancements
- [] Rate limiting
- [] Timeouts and retries for link checking
