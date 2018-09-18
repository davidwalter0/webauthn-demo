package main

import (
	"net/http"

	"github.com/labstack/echo"
)

func Index(c echo.Context) error {
	return c.HTML(http.StatusOK, `<html>
<head>
	<title>WebAuthN demo</title>
	<style>
		.hide {
			display: none;
		}
	</style>
</head>
<body>
	<h1>WebAuthN Demo</h1>

	<p class="hide" id="registerLoading">Registering... Please tap your authenticator.</p>
	<p class="hide" id="loginLoading">Logging in... Please tap your authenticator.</p>

	<form method="post" id="registerForm">
		<input type="text" name="name" id="registerName" />
		<button type="submit">Register</button>
	</form>

	<form method="post" id="loginForm">
		<input type="text" name="name" id="loginName" />
		<button type="submit">Login</button>
	</form>

	<script type="text/javascript">
class WebAuthN {
	static _decodeBuffer(value) {
		return Uint8Array.from(atob(value), c => c.charCodeAt(0));
	}

	static _encodeBuffer(value) {
		return btoa(new Uint8Array(value).reduce((s, byte) => s + String.fromCharCode(byte), ''));
	}

	static _checkStatus(status) {
		return res => {
			if (res.status === status) {
				return res;
			}
			throw new Error(res.statusText);
		};
	}

	register(name) {
		return fetch('/webauthn/registration/start/' + name, {
				method: 'POST'
			})
			.then(WebAuthN._checkStatus(200))
			.then(res => res.json())
			.then(res => {
				res.publicKey.challenge = WebAuthN._decodeBuffer(res.publicKey.challenge);
				res.publicKey.user.id = WebAuthN._decodeBuffer(res.publicKey.user.id);
				if (res.publicKey.excludeCredentials) {
					for (var i = 0; i < res.publicKey.excludeCredentials.length; i++) {
						res.publicKey.excludeCredentials[i].id = WebAuthN._decodeBuffer(res.publicKey.excludeCredentials[i].id);
					}
				}
				return res;
			})
			.then(res => navigator.credentials.create(res))
			.then(credential => {
				return fetch('/webauthn/registration/finish/' + name, {
					method: 'POST',
					headers: {
						'Accept': 'application/json',
						'Content-Type': 'application/json'
					},
					body: JSON.stringify({
						id: credential.id,
						rawId: WebAuthN._encodeBuffer(credential.rawId),
						response: {
							attestationObject: WebAuthN._encodeBuffer(credential.response.attestationObject),
							clientDataJSON: WebAuthN._encodeBuffer(credential.response.clientDataJSON)
						},
						type: credential.type
					}),
				})
			})
			.then(WebAuthN._checkStatus(201))
			.then(res => alert('This authenticator has been registered'))
			.catch(err => {
				console.error(err)
				alert('Failed to register: ' + err);
			});
	}

	login(name) {
		return fetch('/webauthn/login/start/' + name, {
				method: 'POST'
			})
			.then(WebAuthN._checkStatus(200))
			.then(res => res.json())
			.then(res => {
				res.publicKey.challenge = WebAuthN._decodeBuffer(res.publicKey.challenge);
				if (res.publicKey.allowCredentials) {
					for (let i = 0; i < res.publicKey.allowCredentials.length; i++) {
						res.publicKey.allowCredentials[i].id = WebAuthN._decodeBuffer(res.publicKey.allowCredentials[i].id);
					}
				}
				return res;
			})
			.then(res => navigator.credentials.get(res))
			.then(credential => {
				return fetch('/webauthn/login/finish/' + name, {
					method: 'POST',
					headers: {
						'Accept': 'application/json',
						'Content-Type': 'application/json'
					},
					body: JSON.stringify({
						id: credential.id,
						rawId: WebAuthN._encodeBuffer(credential.rawId),
						response: {
							clientDataJSON: WebAuthN._encodeBuffer(credential.response.clientDataJSON),
							authenticatorData: WebAuthN._encodeBuffer(credential.response.authenticatorData),
							signature: WebAuthN._encodeBuffer(credential.response.signature),
							userHandle: WebAuthN._encodeBuffer(credential.response.userHandle),
						},
						type: credential.type
					}),
				})
			})
			.then(WebAuthN._checkStatus(200))
			.then(res => res.json())
			.then(res => alert('You have been logged in to ' + res.name))
			.catch(err => {
				console.error(err)
				alert('Failed to login: ' + err);
			});
	}
}

let registerPending = false;
let loginPending = false;

let w = new WebAuthN();

document.getElementById("registerForm").onsubmit = function(e) {
	e.preventDefault();

	if (registerPending) return;
	registerPending = true;

	document.getElementById("registerLoading").classList.remove("hide");

	const name = document.getElementById("registerName").value;
	w.register(name).then(() => {
		registerPending = false;
		document.getElementById("registerLoading").classList.add("hide");
	});
};

document.getElementById("loginForm").onsubmit = function(e) {
	e.preventDefault();

	if (loginPending) return;
	loginPending = true;

	document.getElementById("loginLoading").classList.remove("hide");

	const name = document.getElementById("loginName").value;
	w.login(name).then(() => {
		loginPending = false;
		document.getElementById("loginLoading").classList.add("hide");
	});
};
	</script>
</body>
</html>`)
}
