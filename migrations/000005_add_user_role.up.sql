CREATE TABLE IF NOT EXISTS public.user_role
(
  user_role_id SERIAL PRIMARY KEY,
  created_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  user_id UUID REFERENCES public.user (user_id) ON DELETE CASCADE NOT NULL,
  role_id INT REFERENCES public.role (role_id) ON DELETE CASCADE NOT NULL
);