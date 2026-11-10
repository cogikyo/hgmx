import { animate, τmid, τmin, τmax } from "../motion.js";

export default function (el) {
  const ribbon = el.querySelector("g");
  animate(ribbon, { scaleY: [1, 1.15, 0.95, 1] }, { duration: τmax });

  const tick = el.querySelector('path[d^="m9 10"]');
  animate(tick, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
  animate(tick, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
