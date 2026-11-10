import { animate, τmax } from "../motion.js";

export default function (el) {
  const helmet = el.querySelector("g");
  animate(helmet, { translateX: [0, 2, -1, 0], rotate: [0, 8, -3, 0] }, { duration: τmax });
}
