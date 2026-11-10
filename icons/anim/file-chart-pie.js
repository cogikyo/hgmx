import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const page = el.querySelector('path[d^="M15.941"]');
  animate(page, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const slice = el.querySelector('path[d^="M9 16"]');
  animate(slice, { translateX: [0, 2, -0.5, 0], translateY: [0, -2, 0.5, 0] }, { duration: τmid, delay: τmin });
}
