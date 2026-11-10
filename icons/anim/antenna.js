import { animate, τmax } from "../motion.js";

export default function (el) {
  const array = el.querySelector("g");
  animate(array, { rotate: [0, -12, 8, 0] }, { duration: τmax });
}
