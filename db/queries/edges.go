package queries

const INSERT_EDGES = `INSERT INTO tasks (destination_id, source_id) VALUES %v RETURNING *;`
const DELETE_EDGES = `DELETE FROM tasks WHERE id in $1;`
