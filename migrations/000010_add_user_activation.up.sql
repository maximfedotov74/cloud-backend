CREATE TABLE IF NOT EXISTS public.user_activation (
  user_activation_id SERIAL PRIMARY KEY,
  activation_account_link UUID DEFAULT NULL,
  end_time timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP + INTERVAL '2 hours',
  user_id UUID UNIQUE REFERENCES public.user (user_id) ON DELETE CASCADE NOT NULL
);