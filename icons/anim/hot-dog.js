import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const bun = el.querySelectorAll('path[d^="M17.1"], path[d^="M6.9"]');
  animate(bun, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const sausage = el.querySelector('path[d^="M21.3"]');
  animate(sausage, { translateX: [0, 1.5, -0.5, 0], translateY: [0, -1.5, 0.5, 0] }, { duration: τmid, delay: τmin });
}
