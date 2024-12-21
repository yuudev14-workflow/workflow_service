package queries

const INSERT_WORKFLOW = `
INSERT INTO workflows (name, description) 
VALUES ($1, $2) 
RETURNING *`

const UPDATE_WORKFLOW = `
UPDATE workflows
SET %v
WHERE id = ($1)
RETURNING *`
