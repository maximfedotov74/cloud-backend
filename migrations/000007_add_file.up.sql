CREATE TABLE IF NOT EXISTS public.file (
  file_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  created_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at timestamp(3),
  title VARCHAR(255) NOT NULL,
  ext VARCHAR(255) NOT NULL,
  size int8 NOT NULL,
  folder_id UUID REFERENCES public.folder (folder_id) ON DELETE CASCADE,
  user_id UUID REFERENCES public.user (user_id) ON DELETE CASCADE NOT NULL
);
