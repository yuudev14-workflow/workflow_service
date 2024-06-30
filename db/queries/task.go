package queries

const UPSERT_TASK = `
INSERT INTO tasks (workflow_id, name, description, parameters)
VALUES $1
ON DUPLICATE KEY UPDATE
	name = VALUES(name),
	description = VALUES(description),
	parameters = VALUES(parameters);
`
const DELETE_TASKS = `DELETE FROM tasks WHERE id in $1`
