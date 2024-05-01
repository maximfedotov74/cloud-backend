package keys

const AuthorizationHeader = "Authorization"
const UserAgentHeader = "User-Agent"

const UserRole = "USER"
const AdminRole = "ADMIN"

type CtxType string

const PsqlUUID = "uuid_generate_v4()"
const PsqlCurrentTimestamp = "CURRENT_TIMESTAMP"

const UserSessionCtx CtxType = "user_session"

const UserTable = "public.user"
const RoleTable = "public.role"
const UserRoleTable = "public.user_role"
const UserActivationTable = "public.user_activation"
const UserChangePasswordCodeTable = "public.change_password_code"
const SessionTable = "public.session"
const FileTable = "public.file"
const FileLinkTable = "public.file_link"
const FolderTable = "public.folder"
const FolderLinkTable = "public.folder_link"

const AccessToken = "access_token"
const RefreshToken = "refresh_token"
