import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const corners = el.querySelector("g[transform-origin]");
  animate(corners, { scale: [1.15, 0.95, 1] }, { duration: τmax });

  const square = el.querySelector("rect");
  animate(square, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
