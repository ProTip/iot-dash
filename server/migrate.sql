CREATE TABLE account(
    id TEXT PRIMARY KEY,
    plan TEXT,
    iot_user_limit INT,
    admin_username TEXT,
    admin_password TEXT,
    count INTEGER DEFAULT 0
);

CREATE TABLE iot_user(
    account_id TEXT,
    user_id TEXT,
    primary key(account_id, user_id)
);

CREATE TABLE bearer_token(
    account_id TEXT,
    token TEXT,
    primary key(account_id, token)
);

INSERT INTO account VALUES(
    "testacct-0000-0000-0000-000000000000",
    "STD",
    5,
    "admin@gmail.com",
    "$2y$12$NgURagEvCLYgdoWwRXgHF.dxTaCkKjlOijgc9j2CmSeXYcjisx7EC",
    0);

INSERT INTO bearer_token VALUES(
    "testacct-0000-0000-0000-000000000000",
    "5b122098ac4800016a2b848a36ab718b3b02ed9e40b6199c20f6cd5cf5ccc3d4");
