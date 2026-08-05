export class LoginRequest {
    constructor(data = {}) {
        this.email = data.email ?? "";
        this.senha = data.senha ?? "";
    }
}