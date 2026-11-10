import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const corners = el.querySelector("g[transform-origin]");
  animate(corners, { scale: [1.15, 0.95, 1] }, { duration: τmax });

  const line = el.querySelector('path[d="M7 12h10"]');
  animate(line, { translateY: [0, -5, 5, 0] }, { duration: τmid, delay: τmin });
}
