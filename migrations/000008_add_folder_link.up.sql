CREATE TABLE IF NOT EXISTS public.folder_link (
  folder_link_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  created_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  validity_period timestamp(3),
  folder_id UUID REFERENCES public.folder (folder_id) ON DELETE CASCADE NOT NULL
);
