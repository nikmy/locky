CREATE OR REPLACE FUNCTION
    set_password(_user_id INT, _service VARCHAR(64), _login VARCHAR(64), _password VARCHAR(64))
    RETURNS VOID
AS
$$
BEGIN
    INSERT INTO users_data (user_id, service, login, password)
    VALUES ($1, $2, $3, $4)
    ON CONFLICT (user_id, service) DO UPDATE SET login       = $3,
                                                 password    = $4,
                                                 last_update = now();
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION
    get_credentials(_user_id INT, _service VARCHAR(64))
    RETURNS TABLE
            (
                _login    VARCHAR(64),
                _password VARCHAR(64)
            )
AS
$$
BEGIN
    RETURN QUERY
        SELECT login, password
        FROM users_data
        WHERE user_id = $1
          AND service = $2;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION
    delete_credentials(_user_id INT, _service VARCHAR(64))
    RETURNS VOID
AS
$$
BEGIN
    DELETE
    FROM users_data
    WHERE user_id = $1
      AND service = $2;
END;
$$ LANGUAGE plpgsql;