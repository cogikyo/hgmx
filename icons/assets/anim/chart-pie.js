import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const pie = el.querySelector('path[d^="M21.21"]');
  animate(pie, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const slice = el.querySelector('path[d^="M21 12"]');
  animate(slice, { translateX: [0, 2.5, -0.5, 0], translateY: [0, -2.5, 0.5, 0] }, { duration: τmid, delay: τmin });
}
