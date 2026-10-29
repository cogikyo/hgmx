import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const badge = el.querySelector('path[d^="M3.85"]');
  animate(badge, { pathLength: [0, 1], opacity: [0.5, 1], rotate: [-180, 0] }, { duration: τmid });
  const x = el.querySelectorAll("line");
  animate(x, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: τmin });
}
