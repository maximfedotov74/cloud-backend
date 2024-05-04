CREATE TABLE IF NOT EXISTS public.user_role
(
  user_role_id SERIAL PRIMARY KEY,
  created_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  user_id UUID REFERENCES public.user (user_id) ON DELETE CASCADE NOT NULL,
  role_id INT REFERENCES public.role (role_id) ON DELETE CASCADE NOT NULL
);


SELECT COUNT(u.user_id)
FROM public.user u
JOIN public.user_role ur1 ON u.user_id = ur1.user_id
JOIN public.role r1 ON ur1.role_id = r1.role_id AND r1.title = 'ADMIN'
JOIN public.user_role ur2 ON u.user_id = ur2.user_id
JOIN public.role r2 ON ur2.role_id = r2.role_id AND r2.title = 'USER';


