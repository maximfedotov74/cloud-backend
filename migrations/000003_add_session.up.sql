CREATE TABLE IF NOT EXISTS public.session
(
  session_id SERIAL PRIMARY KEY,
  refresh_token TEXT NOT NULL UNIQUE,
  created_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  user_agent TEXT NOT NULL,
  ip INET NOT NULL,
  user_id UUID REFERENCES public.user (user_id) ON DELETE CASCADE NOT NULL
);