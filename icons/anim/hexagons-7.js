import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const outline = el.querySelector('path[d^="M5.3 4.3v"]');
  animate(outline, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const walls = el.querySelectorAll('path:not([d^="M5.3 4.3v"])');
  animate(walls, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
