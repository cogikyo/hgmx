import { animate, τmax } from "../motion.js";

export default function (el) {
  const cabin = el.querySelector("g");
  animate(cabin, { rotate: [0, -8, 6, -3, 0] }, { duration: τmax });
}
