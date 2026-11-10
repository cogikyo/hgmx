import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const lens = el.querySelector("g");
  animate(lens, { scale: [1.35, 1] }, { duration: τmax });
  animate(lens, { translateX: [0, -4, 2, 0], translateY: [0, -3, -2, 0] }, { duration: τmax });

  const slash = el.querySelector('path[d^="m13.5"]');
  animate(slash, { scale: [1, 1.25, 1] }, { duration: τmid, delay: τmin });
}
