import { animate, τmax } from "../motion.js";

export default function (el) {
  const bush = el.querySelector("path[transform-origin]");
  animate(bush, { scale: [1, 0.9, 1.08, 1] }, { duration: τmax });
}
