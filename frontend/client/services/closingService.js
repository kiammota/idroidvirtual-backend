import { HttpClient } from "../http.js";

class ClosingServiceClass extends HttpClient {
  async getByDay() {
    return this.get("/closing/dia");
  }

  async getByMonth() {
    return this.get("/closing/mes");
  }
}

/** Singleton */
export const ClosingService = new ClosingServiceClass();
export default ClosingService;
