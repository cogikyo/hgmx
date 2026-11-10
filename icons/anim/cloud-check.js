import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const cloud = el.querySelector('path[d^="M5.516"]');
  animate(cloud, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  const tick = el.querySelector('path[d^="m17 15"]');
  animate(tick, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
  animate(tick, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
