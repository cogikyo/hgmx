import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const page = el.querySelector('path[d^="M4 6.835"]');
  animate(page, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const phones = el.querySelector('path[d^="M2 19"]');
  animate(phones, { scaleX: [1, 0.75, 1.15, 1] }, { duration: τmid, delay: τmin });
}
