import { animate, τmax } from "../motion.js";

export default function (el) {
  const scarf = el.querySelector("g");
  animate(scarf, { rotate: [0, -8, 5, -2, 0] }, { duration: τmax });
}
