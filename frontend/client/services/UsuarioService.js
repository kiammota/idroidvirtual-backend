import { HttpClient } from "../http.js";
import { Usuario } from "../models/UsuarioModel.js";

class UserServiceClass extends HttpClient {
  async create(params) {
    const { empresa_id, email, nome, senha } = params;
    const response = await this.post("/users", { empresa_id, email, nome, senha });
    if (response.success && response.data) {
      response.data = new Usuario(response.data);
    }
    return response;
  }

  async getById(id) {
    const response = await this.get(`/users/${id}`);
    if (response.success && response.data) {
      response.data = new Usuario(response.data);
    }
    return response;
  }

  async elevate(senha) {
    return this.post("/users/elevate", { senha });
  }
}

/** Singleton */
export const UserService = new UserServiceClass();
export default UserService;
