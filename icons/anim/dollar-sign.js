import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const s = el.querySelector('path[d^="M17 5"]');
  animate(s, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const bar = el.querySelector('path[d="M12 3v18"]');
  animate(bar, { scaleY: [0, 1.25, 1] }, { duration: τmid, delay: τmin });
}
