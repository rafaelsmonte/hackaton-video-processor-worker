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
   - The code is sent to **SonarQube** for static analysis, where it is evaluated for:
     - **Potential security vulnerabilities**
     - **Code coverage**
     - **Code duplication**
     - **Dependency checks** (ensuring no critical vulnerabilities exist in third-party libraries)

### PR Approval

- For the PR to be accepted, at least **one reviewer** from the project must approve the change.
- The person who made the commit **cannot be the reviewer**.

## Approval Flow

- When the PR is submitted, it will be reviewed by at least one member of the project.
- If all checks pass (build, tests, SonarQube analysis, and reviewer approval), the code can be merged into the `main` branch.

## CD Process

Once the code is merged into `main`, the **Continuous Deployment (CD)** process is triggered:

### Steps

1. **Upload to ECR**:
   - A **container image** will be built and uploaded to **Amazon Elastic Container Registry (ECR)**.
   - Two versions of the image will be uploaded:
     - **Tag `latest`**: Always pointing to the most recent stable build.
     - **Tag with the commit hash**: Allowing precise version control and easy rollbacks if needed.

2. **Update Container on EKS**:
   - After the image is successfully uploaded to ECR, the container on **Amazon Elastic Kubernetes Service (EKS)** will be updated to reflect the new image.
   - This ensures that the latest version of the code is running in the production environment.

## Considerations

- The `main` branch is kept stable and secure, as all changes go through rigorous testing, SonarQube analysis, and manual reviews before being integrated.
- The **SonarQube** analysis helps prevent security vulnerabilities, maintain high code quality, and ensure that dependencies are safe.
- The **CD process** ensures seamless deployments while maintaining easy rollback options through tagged images.
