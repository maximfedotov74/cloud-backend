package mail

import (
	"fmt"
)

const site = "CloudMax"

func (m *MailService) createActivationTemplate(link string, email string) string {
	l := m.config.AppLink + link
	return fmt.Sprintf(`
  <!DOCTYPE html>
		<html>
			<head>
				<meta charset="UTF-8" />
				<title>Подтверждение регистрации в облачном хранилище %s</title>
				<style>
					body {
						background-color: white;
						color: black;
						font-family: Arial, sans-serif;
					}
					.container {
						max-width: 600px;
						margin: 0 auto;
						padding: 20px;
					}
					h1 {
						color: black;
					}
					.button {
						display: inline-block;
						background-color: yellow;
						color: black;
						padding: 10px 20px;
						text-decoration: none;
						border-radius: 5px;
					}
					.button:hover {
						background-color: #ffd700;
					}
					.footer {
						margin-top: 20px;
						text-align: center;
					}
				</style>
			</head>
			<body>
				<div class="container">
					<h1>Подтверждение регистрации в облачном хранилище %s</h1>
					<p>Здравствуйте, %s!</p>
					<p>
						Спасибо за регистрацию в нашем облачном хранилище. Для завершения регистрации,
						пожалуйста, нажмите на кнопку ниже:
					</p>
					<a href="%s" class="button">Подтвердить регистрацию</a>
					<p>Если у вас возникли вопросы, пожалуйста, свяжитесь с нами.</p>
					<div class="footer">
						<p>С уважением,<br />Команда сайта</p>
					</div>
				</div>
			</body>
		</html>
    `, site, site, email, l)
}

func (m *MailService) createChangePasswordCodeTemplate(code string, email string) string {
	return fmt.Sprintf(`
			<!DOCTYPE html>
			<html>
				<head>
					<meta charset="UTF-8" />
					<title>Смена пароля в облачном хранилище %s</title>
					<style>
						body {
							background-color: white;
							color: black;
							font-family: Arial, sans-serif;
						}
						.container {
							max-width: 600px;
							margin: 0 auto;
							padding: 20px;
						}
						h1 {
							color: black;
						}
						.button {
							display: inline-block;
							background-color: yellow;
							color: black;
							padding: 10px 20px;
							text-decoration: none;
							border-radius: 5px;
						}
						.button:hover {
							background-color: #ffd700;
						}
						.footer {
							margin-top: 20px;
							text-align: center;
						}
					</style>
				</head>
				<body>
					<div class="container">
						<h1>Смена пароля в облачном хранилище %s</h1>
						<p>Здравствуйте, %s!</p>
						<p>Код для смены пароля: %s</p>
						<p>Если у вас возникли вопросы, пожалуйста, свяжитесь с нами.</p>
						<div class="footer">
							<p>С уважением,<br />Команда сайта</p>
						</div>
					</div>
				</body>
			</html>
    `, site, site, email, code)
}
