import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const corners = el.querySelector("g[transform-origin]");
  animate(corners, { scale: [1.15, 0.95, 1] }, { duration: τmax });

  const box = el.querySelector("g:not([transform-origin])");
  animate(box, { translateY: [0, -3, 1, 0] }, { duration: τmid, delay: τmin });
}
