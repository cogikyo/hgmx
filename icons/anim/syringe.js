import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const plunger = el.querySelectorAll('path[d^="m18 2"], path[d^="m17 7"]');
  animate(plunger, { translateX: [0, -3, 1, 0], translateY: [0, 3, -1, 0] }, { duration: τmax });

  const barrel = el.querySelectorAll('path:not([d^="m18 2"]):not([d^="m17 7"])');
  animate(barrel, { translateX: [0, -1, 0], translateY: [0, 1, 0] }, { duration: τmin, delay: τmid });
}
