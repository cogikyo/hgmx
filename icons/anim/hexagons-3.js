import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const outline = el.querySelector('path[d^="M8 4.4"]');
  animate(outline, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const walls = el.querySelectorAll('path[d^="m7 10.5"], path[d="M12 13.5V19"]');
  animate(walls, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
