import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const page = el.querySelector('path[d^="M20 10"]');
  animate(page, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const corners = el.querySelector("g");
  animate(corners, { scale: [1.15, 0.95, 1] }, { duration: τmid, delay: τmin });
}
