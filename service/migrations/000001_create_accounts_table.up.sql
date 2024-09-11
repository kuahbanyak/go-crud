CREATE TABLE accounts (
                          id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                          username VARCHAR(100) UNIQUE,
                          password VARCHAR(100),
                          repeat_password VARCHAR(100),
                          createdat TIMESTAMP,
                          createdby VARCHAR(100),
                          updatedat TIMESTAMP,
                          updatedby VARCHAR(100)
);
CREATE TABLE products (
                          id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                          name VARCHAR(100) UNIQUE,
                          decsription VARCHAR(100),
                          price DECIMAL,
                          stock INT,
                          createdat TIMESTAMP,
                          createdby VARCHAR(100),
                          updatedat TIMESTAMP,
                          updatedby VARCHAR(100)
);