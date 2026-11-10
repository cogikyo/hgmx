import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const corners = el.querySelector("g[transform-origin]");
  animate(corners, { scale: [1.15, 0.95, 1] }, { duration: τmax });

  const lid = el.querySelector('path[d^="M18.944"]');
  animate(lid, { scaleY: [1, 0.1, 1] }, { duration: τmid, delay: τmin });

  const pupil = el.querySelector("circle");
  animate(pupil, { scale: [1, 0, 1.25, 1] }, { duration: τmid, delay: τmin });
}
