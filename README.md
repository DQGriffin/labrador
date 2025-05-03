# Labrador

**Labrador** aims to make serverless infrastructure **simple, repeatable, and manageable**.

It uses **JSON** to define your cloud infrastructure — making it easy to version, review, and deploy serverless projects.

---

## Terminology

| Term | Definition |
| ---- | ---------- |
| **Project** | The application you're building |
| **Resource** | An AWS resource (Lambda, S3, API Gateway, etc.) |
| **Stage** | A group of resources |

---

## Installation

Labrador is distributed as a standalone binary.

1. **Download** the latest release for your system.
2. **Make it executable**:
   ```bash
   chmod +x labrador
   ```
3. **Move it to your PATH** (for example):
   ```bash
   sudo mv labrador /usr/local/bin/labrador
   ```

After installation, you can verify it works by running:

```bash
labrador --version
```

---

## Quickstart

Create an env file with AWS credentials:

```bash
# .env
AWS_ACCESS_KEY_ID=
AWS_SECRET_ACCESS_KEY=
AWS_REGION=
```

Initialize a new Labrador project:

```bash
labrador init --name my_project --env dev --output my_project.json
```

Add a stage to your project:

```bash
labrador add stage --project my_project.json --type s3 --name assets --output buckets.json
```

Inspect your infrastructure:

```bash
labrador inspect --project my_project.json --env-file .env --full
```

Deploy your infrastructure:

```bash
labrador deploy --project my_project.json --env-file .env
```

**More than just S3.**

Labrador can also scaffold and deploy Lambda functions and API Gateways:
```bash
labrador add stage --type lambda --help
labrador add stage --type api --help
```

Each stage type supports customizable configuration, sensible defaults, and full environment interpolation.

**Examples**

See /examples for a sample project deploying several lambda functions, and API gateway with several integrations and routes, and an S3 bucket. You'll need to provide one or more Role ARNs for the lambdas if you choose to deploy it.

---

## AWS Credentials
You have two options for passing AWS credentials to Labrador:

1. Use global flags (--aws-access-key-id, --aws-secret-access-key, --aws-region)
2. Use an environment file defining AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, and AWS_REGION

Note: Labrador will never implicitly read .env. You must use --env-file=.env if you want to use variables defined there.
Labrador will implicitly read .labrador.env though.

---

## Defining Infrastructure

Labrador organizes your infrastructure into **projects** and **stages**, with each configuration defined in a JSON file.

| File | Defines |
| ---- | ---------- |
| **Project config** | Defines your project and settings |
| **Stage config** | Defines one or more cloud resources |

---

### Project Configuration

A **project configuration** defines the overall project metadata, such as name, environment, and variables.

Example:

```json
{
  "name": "my-project",
  "environment": "dev",
  "variables": {
    "config_files": "./stages",
    "version": "1.0"
  },
  "stages": []
}
```

---

## Supported Services

| Service | Create | Update | Delete |
|---------|:------:|:------:|:------:|
| **Lambda** | ✅ | ✅ | ✅ |
| **S3** | ✅ | ✅ | ✅ |
| **API Gateway** | ✅ | ✅ | ✅ |
| **IAM roles** | ✅ | ✅ | ✅ |

---

## Status

Labrador is currently in early development.  
While functional, Labrador is still early in development. Minor rough edges may exist as we continue to refine and expand it.

Feedback, issues, and contributions are welcome!