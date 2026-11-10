import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const axes = el.querySelector('path[d^="M3 3"]');
  animate(axes, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const split = el.querySelector('path[d="M13 3v18"]');
  animate(split, { scaleY: [0, 1] }, { duration: τmid, delay: τmin });

  const cut = el.querySelector('path[d="M3 12h18"]');
  animate(cut, { scaleX: [0, 1] }, { duration: τmid, delay: τmin });
}
