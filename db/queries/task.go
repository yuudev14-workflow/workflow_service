package queries

const GET_TASK_BY_WORKFLOW_ID = `SELECT id, workflow_id, name, description, created_at, updated_at from tasks WHERE workflow_id = $1`
