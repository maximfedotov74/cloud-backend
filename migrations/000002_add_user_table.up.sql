CREATE TABLE IF NOT EXISTS public.user (
  user_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  created_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  email VARCHAR(129) NOT NULL UNIQUE,
  avatar_path TEXT,
  password_hash VARCHAR(255) NOT NULL,
  is_activated boolean NOT NULL DEFAULT false
);

--  id SERIAL PRIMARY KEY