import { animate, τmax } from "../motion.js";

export default function (el) {
  const broom = el.querySelector("g");
  animate(broom, { rotate: [0, -10, 8, -4, 0] }, { duration: τmax });
}
