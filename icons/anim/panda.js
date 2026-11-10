import { animate, τmax } from "../motion.js";

export default function (el) {
  const head = el.querySelector("g");
  animate(head, { rotate: [0, 8, -6, 3, 0], translateY: [0, -1, 0.5, 0, 0] }, { duration: τmax });
}
