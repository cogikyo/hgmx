import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const note = el.querySelector("g");
  animate(note, { rotate: [0, -12, 4, 0], translateY: [0, -4, 1, 0] }, { duration: τmax });

  const fold = el.querySelector('path[d^="M15 3v5"]');
  animate(fold, { scale: [0, 1.25, 1] }, { duration: τmin, delay: τmid });

  const tick = el.querySelector('path[d^="m15 19"]');
  animate(tick, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
