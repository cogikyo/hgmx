import { animate, τmax } from "../motion.js";

export default function (el) {
  const reader = el.querySelector("g");
  animate(reader, { translateY: [0, 1.5, -0.5, 0] }, { duration: τmax });
}
