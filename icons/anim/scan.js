import { animate, τmax } from "../motion.js";

export default function (el) {
  const corners = el.querySelector("g[transform-origin]");
  animate(corners, { scale: [1.15, 0.95, 1] }, { duration: τmax });
}
