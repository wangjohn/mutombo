CREATE TABLE IF NOT EXISTS requests(
  id bigserial PRIMARY KEY,
  blocking boolean NOT NULL,
  method text NOT NULL,
  url text NOT NULL,
  finished boolean NOT NULL DEFAULT FALSE,
  body text,
  status_code integer
);

CREATE TABLE IF NOT EXISTS headers(
  id bigserial PRIMARY KEY,
  request_id bigserial references requests(id),
  name text NOT NULL,
  value text NOT NULL
);
CREATE INDEX headers_request_id ON headers(request_id);
