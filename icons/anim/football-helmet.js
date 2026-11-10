import { animate, τmax } from "../motion.js";

export default function (el) {
  const helmet = el.querySelector("g");
  animate(helmet, { translateY: [0, -2, 0], rotate: [0, -12, 4, 0] }, { duration: τmax });
}
