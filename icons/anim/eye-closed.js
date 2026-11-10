import { animate, τmax } from "../motion.js";

export default function (el) {
  const lashes = el.querySelectorAll("path[transform-origin]");
  animate(lashes, { scale: [1, 0.4, 1.25, 1] }, { duration: τmax });
}
