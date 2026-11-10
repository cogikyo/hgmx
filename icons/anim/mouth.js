import { animate, τmax } from "../motion.js";

export default function (el) {
  const mouth = el.querySelector("g");
  animate(mouth, { scaleY: [1, 1.25, 0.9, 1] }, { duration: τmax });
}
