CREATE TABLE list (
  id bigserial primary key not null,
  name varchar(100) not null,
  description varchar(255) not null default '',
  pins integer not null default 0,
  public boolean not null default false
);

CREATE TABLE "user" (
  id bigserial primary key not null,
  status integer not null default 1,
  username varchar(60) not null,
  email varchar(100) not null,
  password varchar(255) not null,
  createdat TIMESTAMP not null default CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX user_email_idx ON "user" (email);
CREATE UNIQUE INDEX user_username_idx ON "user" (username);

CREATE TABLE user_has_list (
  id bigserial primary key not null,
  role integer not null default 1,
  list_id bigint not null,
  user_id bigint not null,
  CONSTRAINT fk_user_user_id FOREIGN KEY (user_id) REFERENCES "user"(id),
  CONSTRAINT fk_list_list_id FOREIGN KEY (list_id) REFERENCES list(id)
);

CREATE UNIQUE INDEX user_has_list_list_id_user_id_idx ON user_has_list (list_id, user_id);

CREATE TABLE token (
  id bigserial primary key not null,
  hash varchar(255) not null,
  until TIMESTAMP not null,
  created_at TIMESTAMP not null,
  user_id bigint not null,
  CONSTRAINT fk_user_user_id FOREIGN KEY (user_id) REFERENCES "user"(id)
);

CREATE TABLE pin (
  id bigserial primary key not null,
  title varchar(255) not null,
  url varchar(255) not null,
  creator_id bigint not null,
  list_id bigint not null default 0,
  created_at TIMESTAMP not null default CURRENT_TIMESTAMP,
  CONSTRAINT fk_user_creator_id FOREIGN KEY (creator_id) REFERENCES "user"(id),
  CONSTRAINT fk_list_list_id FOREIGN KEY (list_id) REFERENCES list(id)
);

CREATE TABLE tag (
  id bigserial primary key not null,
  pin_id bigint not null,
  name varchar(30) not null,
  CONSTRAINT fk_pin_pin_id FOREIGN KEY (pin_id) REFERENCES pin(id)
);

CREATE UNIQUE INDEX tag_pin_id_name_idx ON tag (pin_id, name);
