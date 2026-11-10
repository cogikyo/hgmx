import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const key = el.querySelector('path[d^="M3 3"]');
  animate(key, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const alt = el.querySelector('path[d="M14 3h7"]');
  animate(alt, { scaleX: [0, 1.25, 1] }, { duration: τmid, delay: τmin });
}
