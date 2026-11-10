import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const clock = el.querySelector("g");
  animate(clock, { rotate: [0, -12, 12, -8, 8, 0] }, { duration: τmax });

  const tick = el.querySelector('path[d^="m9 13"]');
  animate(tick, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
  animate(tick, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
