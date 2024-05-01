package msg

const (
	UserNotFound                   = "Пользователь не найден!"
	UserExists                     = "Пользователь с таким email уже существует!"
	UserActivationError            = "Произошла ошибка при активации! Попробуйте в другой раз."
	UserActivationNotFound         = "Пользователя с такой ссылкой не существует!"
	UserErrorWhenAddActivationLink = "Произошла ошибка при создании ссылки для активации аккаунта!"
	UserActionLinkAlreadyExists    = "Срок действия текущей ссылки еще не истек! Вы можете активировать акканут по ней."
	UpdatePasswordError            = "Ошибка при смене пароля!"
	BadPassword                    = "Введенный пароль не совпадает с текущим!"
	BadNewPassword                 = "Новый пароль должен отличаться от старого!"
	ChangePasswordCodeNotFound     = "Неверный код или вышел его срок действия!"
	ChangePasswrodError            = "Ошибка при смене пароля!"
	CreateChangeCodeError          = "Ошибка при создании кода для смены пароля!"
	UpdateUserError                = "Ошибка при обновлении информации пользователя!"
)
