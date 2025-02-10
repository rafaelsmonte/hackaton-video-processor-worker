# CI/CD Workflow - Application

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

## CD Process

Once the code is merged into `main`, the **Continuous Deployment (CD)** process is triggered:

### Steps

1. **Upload to ECR**:
   - A **container image** will be built and uploaded to **Amazon Elastic Container Registry (ECR)**.
   - Two versions of the image will be uploaded:
     - **Tag `latest`**: This will always point to the most recent stable build.
     - **Tag with the commit hash**: This will allow for precise version control and easy rollbacks to any previous version of the application.

2. **Update Container on EKS**:
   - After the image is successfully uploaded to ECR, the container on **Amazon Elastic Kubernetes Service (EKS)** will be updated to reflect the new image.
   - This will ensure that the latest version of the code is running in the production environment.

## Considerations

- The `main` branch will always be kept stable and secure, as changes go through rigorous testing and reviews before being integrated.
- The CD process ensures that the latest code is continuously deployed to the production environment via ECR and EKS.
- Uploading both the `latest` tag and the commit hash tag allows for easy version tracking and quick rollbacks if needed.
