import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const body = el.querySelector('path[d^="M7 20"]');
  animate(body, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const bar = el.querySelector('path[d="M4 16h6"]');
  animate(bar, { scaleX: [0, 1.25, 1] }, { duration: τmid, delay: τmin });
}
