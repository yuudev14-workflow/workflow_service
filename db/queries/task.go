package queries

const GET_TASK_BY_WORKFLOW_ID = `SELECT id, workflow_id, name, description, created_at, updated_at from tasks WHERE workflow_id = $1`

const UPSERT_TASK = `
INSERT INTO tasks (workflow_id, name, description)
VALUES %v
ON CONFLICT (unique_workflow_name) DO UPDATE
	description = VALUES(description),
	updated_at = NOW();
`
const DELETE_TASKS = `DELETE FROM tasks WHERE id in $1;`
