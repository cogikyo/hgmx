import { animate, τmax } from "../motion.js";

export default function (el) {
  const cart = el.querySelector("g");
  animate(cart, { translateX: [0, 4, -2, 0], rotate: [0, -8, 4, 0] }, { duration: τmax });
}
