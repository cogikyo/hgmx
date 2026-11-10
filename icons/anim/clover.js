import { animate, τmax } from "../motion.js";

export default function (el) {
  const clover = el.querySelector("g");
  animate(clover, { rotate: [0, -12, 8, -4, 0] }, { duration: τmax });
}
