import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const badge = el.querySelector('path[d^="M3.85"]');
  animate(badge, { pathLength: [0, 1], opacity: [0.5, 1], rotate: [-180, 0] }, { duration: τmid });
  const tick = el.querySelector('path[d^="m16 9"]');
  animate(tick, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
  animate(tick, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
