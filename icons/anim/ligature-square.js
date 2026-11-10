import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector("rect");
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const f = el.querySelectorAll('path[d^="M14.4 8"], path[d="M7 12h3"], path[d="M10 17H7"]');
  animate(f, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
  const i = el.querySelectorAll('path[d="M14 12h1.5v5"], path[d="M17 17h-3"]');
  animate(i, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: τmid });
}
