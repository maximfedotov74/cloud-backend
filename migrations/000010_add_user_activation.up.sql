CREATE TABLE IF NOT EXISTS public.user_activation (
  user_activation_id SERIAL PRIMARY KEY,
  activation_account_link UUID UNIQUE DEFAULT uuid_generate_v4(),
  end_time timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP + INTERVAL '30 minutes',
  user_id UUID UNIQUE REFERENCES public.user (user_id) ON DELETE CASCADE NOT NULL
);