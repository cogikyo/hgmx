import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const halves = el.querySelectorAll('path:not([d="M2 12h20"])');
  animate(halves, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const cut = el.querySelector('path[d="M2 12h20"]');
  animate(cut, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
