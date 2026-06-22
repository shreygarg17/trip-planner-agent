# Workspace Design Guidelines

Always adhere to the following architectural guidelines for this workspace:

1. **Root-Level Package Structure**:
   - `main.go` resides at the root level.
   - All modules (`models`, `service`, `repo`, `clients`, `config`, `utils`) reside as direct subdirectories of the root workspace.
   - `cmd/api/` contains all API HTTP route handler code.

2. **Interface-Led Design & Package Setup**:
   - Every core module defines its contracts in an `interfaces/` subpackage (e.g., `service/interfaces/`, `clients/interfaces/`, `repo/interfaces/`).
   - Implementations are placed directly under the module folder (e.g., `service/planner.go`).
   - The `models` package contains only data structures (requests, responses, and entities) in files like `requests.go`, `responses.go`, and `entities.go`.

3. **Method Constraints**:
   - Every method should be at max 50 lines. Large functions must be split into concise helper sub-methods.

4. **Configuration & Utilities**:
   - Configuration resides in `config/` (specifically `config.go` for definitions and `env.go` for environment-based loading).
   - Reusable generic utilities reside in `utils/utils.go`.

5. **Test Organization**:
   - Keep service-level unit tests in a dedicated folder inside `service/` (specifically `service/tests/`).

6. **Validations & Constants**:
   - Validation logic (such as parameter bounds and presence checks) resides in `validations/` under root (with interfaces in `validations/interfaces/`).
   - Constant definitions reside in `constants/` under root, using separate files per type (e.g. `error_constants.go` for error string values).

