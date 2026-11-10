import { animate, τmax } from "../motion.js";

export default function (el) {
  const trousers = el.querySelector("g");
  animate(trousers, { rotate: [0, -8, 5, -2, 0] }, { duration: τmax });
}
