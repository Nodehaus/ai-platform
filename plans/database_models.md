# Data models for the AI platform

This is a living document to write down ideas and planning for the data model of our AI platform app. We want users to
allow to create training data and then finetune an LLM. All models have `created_at` and `updated_at` fields. All IDs
are UUIDs. We already have a `User` model in the app, so this is not in the current plan. Foreign keys are not in the
models, we will add those during implementation.

## Project

The main entity of the app will be `Project`. A project owns zero or more `TrainingDataset` and `Finetune`.

### Model sketch

-   type Project
    -   name: string (required)
    -   owner: User (1, required)
    -   training_dataset: TrainingDataset (latest version, 0 or 1)
    -   finetune: Finetune (latest version, 0 or 1)
    -   status: enum of [ACTIVE, ARCHIVED, DELETED] (required)

## Corpus

A corpus has a unique name and an S3 path. The path contains a list of JSON files that we will use as input to generate
the training data.

### Model sketch

-   type Corpus
    -   name: string (required)
    -   s3_path: string (required)

## TrainingDataset

A project will own zero or more items of `TrainingDataset`. The training dataset is versioned, with incremental version
numbers.

### Model sketch

-   type TrainingDataset
    -   version: int (required)
    -   generate_model: string
    -   generate_model_runner: string
    -   generate_gpu_info_card: string
    -   generate_gpu_info_total_gb: float (rounded to 2 decimals)
    -   generate_gpu_info_cuda_version: string
    -   input_field: string
    -   output_field: string
    -   total_generation_time_seconds: float (rounded to 2 decimals)
    -   generate_prompt_history: list of Prompt (history of all prompts that the user saved, except the current_prompt)
    -   generate_prompt: Prompt (required)
    -   corpus: Corpus (required)
    -   language_iso: string (3-letter ISO code, required)
    -   status: enum of [PLANNING, RUNNING, ABORTED, FAILED, DONE, DELETED] (required)
    -   field_names: list of string (required)
    -   data: list of TrainingDataItem

Each `TrainingDataItem` is one example for training, validation and/or evaluation:

-   type TrainingDataItem
    -   values: list of string (required)
    -   corrects: ID of TrainingDataItem
    -   source_document: string
    -   source_document_start: string
    -   source_document_end: string
    -   generation_time_seconds: float (rounded to 2 decimals, required)
    -   deleted: boolean (required, default False)

The list of values are the same length and order as the `field_names` in `TrainingDataset`. When the user edits one
`TrainingDataItem` we add a new database entry where the `corrects` field points to the original `TrainingDataItem`.
Then we replace the item in the `data` field of the `TrainingDataset`.

## Prompt

A `Prompt` is a string with a version.

### Model sketch

-   type Prompt
    -   version: int (required)
    -   text: string (required)

## Finetune

The `Finetune` stored information about the model training and the final model.

### Model sketch

-   type Finetune
    -   version: int (required)
    -   model_name: string (required)
    -   base_model_name: string (required)
    -   model_size_gb: int
    -   model_size_parameter: int
    -   model_dtype: string
    -   model_quantization: string
    -   inference_samples: list of InferenceSample
    -   training_dataset: TrainingDataset (required)
    -   training_time_seconds: float (rounded to 2 decimals)
    -   status: enum of [RUNNING, ABORTED, FAILED, DONE, DELETED] (required)

The `InferenceSample` contains generated output with their input from the validation dataset, we create those during training at specific training steps:

-   type InferenceSample
    -   at_step: int
    -   items: list of [input: string, output: string] pairs

## Status Transitions

### Project Status

```
ACTIVE → ARCHIVED (user archives project)
ACTIVE → DELETED (user deletes project)
ARCHIVED → ACTIVE (user reactivates project)
ARCHIVED → DELETED (user deletes archived project)
```

### TrainingDataset Status

```
PLANNING → RUNNING (generation starts)
PLANNING → DELETED (user cancels before generation)
RUNNING → DONE (generation completes successfully)
RUNNING → FAILED (generation encounters error)
RUNNING → ABORTED (user cancels during generation)
FAILED → PLANNING (user retries after fixing issues)
ABORTED → PLANNING (user restarts generation)
DONE → DELETED (user removes dataset)
```

### Finetune Status

```
RUNNING → DONE (training completes successfully)
RUNNING → FAILED (training encounters error)
RUNNING → ABORTED (user cancels training)
FAILED → RUNNING (user retries after fixing issues)
ABORTED → RUNNING (user restarts training)
DONE → DELETED (user removes model)
```

Note: DELETED status is typically a soft delete - the record remains in database but is hidden from user interface.

## Versioning Implementation

### TrainingDataset Versioning

```
- Version numbers are auto-incremented integers starting from 1
- Each project maintains its own version sequence
- New version is created when:
  * User changes generation parameters (model, prompt, corpus etc.) and re-runs dataset generation
  * User reruns generation after FAILED/ABORTED

Database approach:
- Single table with composite key (project_id, version)
- Query latest: SELECT * FROM training_datasets WHERE project_id = ? ORDER BY version DESC LIMIT 1
- Project.training_dataset references the latest version
```

### Finetune Versioning

```
- Version numbers are auto-incremented integers starting from 1
- Each project maintains its own version sequence
- New version is created when:
  * User changes training parameters
  * User selects different training dataset
  * User reruns training after FAILED/ABORTED

Database approach:
- Single table with composite key (project_id, version)
- Query latest: SELECT * FROM finetunes WHERE project_id = ? ORDER BY version DESC LIMIT 1
- Project.finetune references the latest version
```
