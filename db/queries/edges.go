package queries

const INSERT_EDGES = `INSERT INTO tasks (destination_id, source_id) VALUES $1 RETURNING *;`
const DELETE_EDGES = `DELETE FROM tasks WHERE id in $1;`
