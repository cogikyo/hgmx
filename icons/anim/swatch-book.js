import { animate, τmax } from "../motion.js";

export default function (el) {
  const book = el.querySelector("g");
  animate(book, { rotate: [0, -6, 3, 0] }, { duration: τmax });
}
