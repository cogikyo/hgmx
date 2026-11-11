import { animate, τmax } from "../motion.js";

export default function (el) {
  const reader = el.querySelector("g");
  animate(reader, { translateY: [0, 3, -1, 0] }, { duration: τmax });
}
