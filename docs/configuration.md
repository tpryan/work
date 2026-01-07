# User Configuration Format

The `work` tool uses a YAML configuration file to define data sources, classification rules, and reporting destinations for a specific user. This file is typically located in the `users/` directory (e.g., `users/username.yaml`).

## Top-Level Fields

| Field | Type | Description |
| :--- | :--- | :--- |
| `spread_sheet_id` | String | The ID of the Google Sheet where the report will be generated. |
| `github_user` | String | The GitHub username to query for issues and pull requests. |
| `query_drive` | Boolean | If `true`, the tool will search Google Drive for files owned by the user. |
| `sources` | List | A list of source identifiers. These typically correspond to raw data sheets or internal collectors (e.g., `"Source - Github"`, `"Source - DriveFiles"`). |
| `destinations` | List | Defines the output tabs in the Google Sheet and the data filters for each. |
| `classifiers` | Object | Rules for assigning `Project` and `Subproject` labels to artifacts based on their properties. |

## Destinations

Each item in the `destinations` list defines a report to be generated in a specific tab of the Google Sheet.

| Field | Type | Description |
| :--- | :--- | :--- |
| `sheet` | String | The name of the Google Sheet tab to write to. |
| `sort` | String | (Optional) Sorting method. Use `"report"` for a structured grouping by Project > Subproject > Type. Default is chronological. |
| `summary` | Boolean | (Optional) If `true`, a summary markdown report may be generated or appended. |
| `criteria` | Object | Filters used to select artifacts for this destination. |

### Criteria

| Field | Type | Description |
| :--- | :--- | :--- |
| `start` | Date | Include artifacts created or updated on or after this date (Format: `YYYY-MM-DD`). |
| `end` | Date | Include artifacts created or updated before this date (Format: `YYYY-MM-DD`). |
| `project` | String | (Optional) Only include artifacts matching this project name (case-insensitive). |

## Classifiers

Classifiers are used to automatically tag artifacts with a `Project` and `Subproject`. The tool processes artifacts against these rules; the first matching rule wins (or stamps the artifact).

### Exclusions

| Field | Type | Description |
| :--- | :--- | :--- |
| `exclusions` | List | A list of specific URLs to completely ignore during processing. |

### Lists

The `lists` section defines specific rules for classification.

| Field | Type | Description |
| :--- | :--- | :--- |
| `project` | String | The project name to assign if a match is found. |
| `subproject` | String | (Optional) The subproject name to assign. |
| `links` | List | (Optional) A list of exact URLs that explicitly belong to this project/subproject. |
| `contains` | Object | (Optional) Rules for substring matching. |

#### Contains

| Field | Type | Description |
| :--- | :--- | :--- |
| `title` | List | Match if the artifact's title contains any of these strings (case-insensitive). |
| `link` | List | Match if the artifact's URL contains any of these strings (case-insensitive). |

## Example Configuration

```yaml
spread_sheet_id: "1BgQ2c6QRkhlGN2Pg7tVqrNUxAs5b66f0R29csePYMv0"
github_user: "octocat"
query_drive: true

sources:
  - "Source - Github"
  - "Source - DriveFiles"

destinations:
  - sheet: "2024 Annual"
    criteria:
      start: 2024-01-01
      end: 2025-01-01
  - sheet: "Project Alpha"
    sort: "report"
    criteria:
      project: "Alpha"
      start: 2024-01-01
      end: 2025-01-01

classifiers:
  exclusions:
    - "https://docs.google.com/ignore-me"
  lists:
    - project: "Alpha"
      subproject: "Phase 1"
      links:
        - "https://github.com/org/repo/issues/123"
      contains:
        title:
          - "alpha"
          - "prototype"
        link:
          - "repo-alpha"
```
