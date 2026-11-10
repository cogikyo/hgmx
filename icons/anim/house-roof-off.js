import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const walls = el.querySelector('path[d^="M4 14"]');
  animate(walls, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const roof = el.querySelector("g");
  animate(roof, { scaleY: [1, 0.9, 1] }, { duration: τmid });

  const slash = el.querySelector('path[d="m2 2 20 20"]');
  animate(slash, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
