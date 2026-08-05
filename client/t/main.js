import { API_BASE_URL } from "../env.js";
import { ProductService } from "../services/productService.js";

function main() {
    console.log("hello!");

    const productService = new ProductService(API_BASE_URL);

    console.log(productService);
}

main();