import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const lens = el.querySelector("g");
  animate(lens, { scale: [1.35, 1] }, { duration: τmax });
  animate(lens, { translateX: [0, -4, 2, 0], translateY: [0, -3, -2, 0] }, { duration: τmax });

  const tick = el.querySelector('path[d^="m8 11"]');
  animate(tick, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
  animate(tick, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
