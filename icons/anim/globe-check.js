import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const globe = el.querySelector('path[d^="M2 12h20"]');
  animate(globe, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const tick = el.querySelector('path[d^="m15 6 "]');
  animate(tick, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
  animate(tick, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
