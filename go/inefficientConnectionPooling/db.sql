-- employee table
CREATE TABLE IF NOT EXISTS employee (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL
);

-- inserts 100 employee records
DO $$
BEGIN
    FOR i IN 1..100 LOOP
        INSERT INTO employee (name, email)
        VALUES (
            'Employee' || i,
            'employee' || i || '@example.com'
        );
    END LOOP;
END $$;

