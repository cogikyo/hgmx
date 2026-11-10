import { animate, τmax } from "../motion.js";

export default function (el) {
  const glass = el.querySelector("g");
  animate(glass, { translateX: [0, 1.5, -0.5, 0] }, { duration: τmax });
}
