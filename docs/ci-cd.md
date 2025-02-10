# CI/CD Workflow - Application

## Development Flow

- **`main` Branch**: The `main` branch is protected, meaning direct commits to it are not allowed.
- **Method for Pushing Code to `main`**: The only way to push code to the `main` branch is through **Pull Requests (PRs)**.

## PR Requirements

For a Pull Request to be accepted into the `main` branch, the following verification and validation process must be followed:

### Mandatory Steps

1. **Build**: The code in the PR must pass the build step successfully.
2. **Tests & Analysis**:
   - The code must be tested, and the **test coverage must be greater than 80%** for the PR to be accepted.
   - The code is sent to **SonarQube**, which performs:
     - Security vulnerability analysis.
     - Code coverage validation.
     - Code duplication detection.
     - Dependency analysis to check for vulnerabilities.
3. **PR Approval**:
   - For the PR to be accepted, at least **one reviewer** from the project must approve the change.
   - The person who made the commit **cannot be the reviewer**.

## Approval Flow

- When the PR is submitted, it will be reviewed by at least one member of the project.
- After the review and approval by the reviewer, the code can be merged into the `main` branch.

## CD Process

Once the code is merged into `main`, the **Continuous Deployment (CD)** process is triggered:

### Steps

1. **Build & Tests**:
   - The application is built again to ensure consistency.
   - The same **test suite** used in CI runs again to confirm that the merged code remains stable.
   - The code is analyzed again by **SonarQube**, verifying:
     - Security vulnerabilities.
     - Code coverage.
     - Code duplication.
     - Dependencies.

2. **Upload to ECR**:
   - A **container image** is built and uploaded to **Amazon Elastic Container Registry (ECR)**.
   - Two versions of the image are pushed:
     - **Tag `latest`**: Always points to the most recent stable build.
     - **Tag with the commit hash**: Enables precise version control and easier rollbacks.

3. **Update Container on EKS**:
   - After the image is uploaded to ECR, the container on **Amazon Elastic Kubernetes Service (EKS)** is updated.
   - This ensures the latest validated version is running in the production environment.

## Considerations

- The `main` branch remains stable and secure, as changes go through rigorous testing and reviews before integration.
- The **CD process also includes testing and SonarQube analysis**, ensuring that any new deployment maintains high quality and security.
- Uploading both the `latest` tag and the commit hash tag allows for easy version tracking and quick rollbacks if needed.
