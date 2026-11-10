import { animate, τmax } from "../motion.js";

export default function (el) {
  const helmet = el.querySelectorAll("path, circle");
  animate(helmet, { translateY: [0, -2, 0.5, 0] }, { duration: τmax });
}
