import { animate, τmax } from "../motion.js";

export default function (el) {
  const guitar = el.querySelector("g");
  animate(guitar, { rotate: [0, -12, 10, -6, 0] }, { duration: τmax });
}
