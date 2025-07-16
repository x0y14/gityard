-- using mysql

create table users (
    id bigint not null auto_increment,
    email varchar(255), -- 退会時に解放のためnull許容
    is_deleted tinyint(1) not null default 0, -- 0=有効、1=退会済み
    created_at datetime default current_timestamp,
    updated_at datetime default current_timestamp on update current_timestamp,

    primary key(id),
    unique index uq_idx_users_email (email)
);
create table user_credentials (
    user_id bigint not null,
    hashed_password varchar(255) not null,
    created_at datetime default current_timestamp,
    updated_at datetime default current_timestamp on update current_timestamp,

    primary key(user_id),
    foreign key(user_id) references users(id) on delete restrict
);
create table user_refresh_tokens (
    user_id bigint not null,
    hashed_refresh_token varchar(255) not null,
    expires_at datetime not null, -- 定期的にDBスキャンして期限切れを削除するため
    created_at datetime default current_timestamp,
    updated_at datetime default current_timestamp on update current_timestamp,

    primary key(user_id),
    unique index uq_idx_users_hashed_refresh_token (hashed_refresh_token),
    foreign key(user_id) references users(id) on delete cascade -- ユーザ削除時に一緒に消す
);

create table handlenames (
    id bigint not null auto_increment,
    handlename varchar(255) not null,
    created_at datetime default current_timestamp,

    primary key(id),
    unique index uq_idx_handlenames_handlename (handlename)
);

create table accounts (
    id bigint not null auto_increment,
    user_id bigint not null,
    handlename_id bigint,
    kind smallint not null default 1, -- 1=個人, 2=組織
    is_deleted tinyint(1) not null default 0, -- 0=有効、1=退会済み
    created_at datetime default current_timestamp,
    updated_at datetime default current_timestamp on update current_timestamp,

    primary key(id),
    unique index uq_idx_accounts_handlename_id (handlename_id),
    foreign key(user_id) references users(id) on delete restrict,
    foreign key(handlename_id) references handlenames(id) on delete restrict
);
create table account_publickeys ( -- openssh format
    id bigint not null auto_increment,
    account_id bigint not null,
    fulltext text not null,
    algorithm varchar(50) not null,
    keybody text not null,
    comment varchar(255) not null,
    fingerprint varchar(255) not null,
    created_at datetime default current_timestamp,

    primary key(id),
    unique (account_id, fingerprint),
    index idx_account_publickeys_fingerprint (fingerprint), -- for account's pubkey list
    foreign key(account_id) references accounts(id) on delete cascade, -- account削除時に一緒に消す
);
create table account_profiles (
    account_id bigint not null,
    displayname varchar(255) not null default "unknown",
    iconpath varchar(255) not null default "noimage001",
    is_private tinyint(1) not null default 0, -- 0=公開, 1=非公開
    created_at datetime default current_timestamp,
    updated_at datetime default current_timestamp on update current_timestamp,

    primary key(account_id),
    foreign key(account_id) references accounts(id) on delete cascade, -- account削除時に一緒に消す
    index idx_account_profiles_displayname (displayname)
);


create table repositories (
    id bigint not null auto_increment,
    owner_account_id bigint,
    name varchar(255) not null,
    is_private tinyint(1) not null default 0, -- 0=公開, 1=非公開
    created_at datetime default current_timestamp,
    updated_at datetime default current_timestamp on update current_timestamp,

    primary key(id),
    foreign key(owner_account_id) references accounts(id) on delete set null, -- on delete, dont foget set is_private=1
    unique index uq_idx_repositoris_owner_name (owner_account_id, name) -- disallow same name per account
);