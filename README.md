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

Initialize a new Labrador project:

```bash
labrador init
```

This will give you a basic project configuration. Feel free to rename the file to anything you like. 
Future releases will support scaffolding out stage configurations as well, but for now, just take a look at
the samples in /templates.

Inspect your infrastructure:

```bash
labrador --project=project.json --env-file=.env inspect (optionally --verbose)
```

Deploy your infrastructure:

```bash
labrador --project=project.json --env-file=.env deploy
```

---

## AWS Credentials
You have two options for passing AWS credentials to Labrador:

1. Use flags (--aws-access-key-id, --aws-secret-access-key, --aws-region)
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
| **API Gateway** | ✅ | ❌ | ✅ |

---

## Status

Labrador is currently in early development (`v0.1.0`).  
While functional, Labrador is still early in development. Minor rough edges may exist as we continue to refine and expand it.

Feedback, issues, and contributions are welcome!