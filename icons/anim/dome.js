import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const base = el.querySelectorAll('path[d^="M22 19"], path[d="M6 12v9"], path[d="M18 12v9"]');
  animate(base, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const dome = el.querySelector("g");
  animate(dome, { scaleY: [1, 0.8, 1.05, 1] }, { duration: τmid, delay: τmin });
}
