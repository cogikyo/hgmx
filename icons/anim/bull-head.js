import { animate, τmax } from "../motion.js";

export default function (el) {
  const head = el.querySelector("g");
  animate(head, { rotate: [0, -10, 8, -4, 0] }, { duration: τmax });
}
