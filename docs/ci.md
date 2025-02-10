# CI Workflow - Application

## Development Flow

- **`main` Branch**: The `main` branch is protected, meaning direct commits to it are not allowed.
- **Method for Pushing Code to `main`**: The only way to push code to the `main` branch is through **Pull Requests (PRs)**.

## PR Requirements

For a Pull Request to be accepted into the `main` branch, the following verification and validation process must be followed:

### Mandatory Steps

1. **Build**: The code in the PR must pass the build step successfully.
2. **Tests**: The code added must be tested. The test coverage must be greater than **80%** for the PR to be accepted.
3. **Test Coverage**: The test coverage should be monitored and must be above **80%** to ensure the quality of the inserted code.

### PR Approval

- For the PR to be accepted, at least **one reviewer** from the project must approve the change.
- The person who made the commit **cannot be the reviewer**.

## Approval Flow

- When the PR is submitted, it will be reviewed by at least one member of the project.
- After the review and approval by the reviewer, the code can be merged into the `main` branch.

## Considerations

- The `main` branch will always be kept stable and secure, as changes go through rigorous testing and reviews before being integrated.
