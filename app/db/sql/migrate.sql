CREATE OR REPLACE FUNCTION
    set_password(user_id INT, service VARCHAR(64), login VARCHAR(64), password VARCHAR(64))
    RETURNS VOID
AS
$$
BEGIN
    UPDATE users_data
    SET login       = $3,
        password    = $4,
        last_update = now()
    WHERE user_id = $1
      AND service = $2;
    INSERT INTO users_data (user_id, service, login, password)
    VALUES ($1, $2, $3, $4)
    ON CONFLICT DO NOTHING;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION
    get_credentials(user_id INT, service VARCHAR(64))
    RETURNS TABLE
            (
                login    VARCHAR(64),
                password VARCHAR(64)
            )
AS
$$
BEGIN
    RETURN QUERY
        SELECT login, password
        FROM users_data
        WHERE users_data.user_id = $1
          AND users_data.service = $2;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION
    delete_credentials(user_id INT, service VARCHAR(64))
    RETURNS VOID
AS
$$
BEGIN
    DELETE
    FROM users_data
    WHERE users_data.user_id = $1
      AND users_data.service = $2;
END;
$$ LANGUAGE plpgsql;