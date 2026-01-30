CREATE TABLE history (
    uuid UUID PRIMARY KEY,
    book_title VARCHAR(255),
    user_request JSONB NOT NULL ,

    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    error TEXT,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    update_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    chapters JSONB
);