import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const corners = el.querySelector("g[transform-origin]");
  animate(corners, { scale: [1.15, 0.95, 1] }, { duration: τmax });

  const mouth = el.querySelector('path[d^="M8 14"]');
  animate(mouth, { scaleY: [0, 1.15, 1] }, { duration: τmid, delay: τmin });

  const eyes = el.querySelectorAll('path[d$="9h.01"]');
  animate(eyes, { translateY: [1, -1, 0] }, { duration: τmin, delay: τmid });
}
