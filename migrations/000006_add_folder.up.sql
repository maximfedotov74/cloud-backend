CREATE TABLE IF NOT EXISTS public.folder (
  folder_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  created_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at timestamp(3),
  title VARCHAR(255) NOT NULL,
  parent_folder_id UUID REFERENCES public.folder (folder_id) ON DELETE CASCADE CHECK (folder_id != parent_folder_id),
  user_id UUID REFERENCES public.user (user_id) ON DELETE CASCADE NOT NULL
);
