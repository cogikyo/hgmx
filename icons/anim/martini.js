import { animate, τmax } from "../motion.js";

export default function (el) {
  const glass = el.querySelector("g");
  animate(glass, { rotate: [0, -10, 6, 0] }, { duration: τmax });
}
