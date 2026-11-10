import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const dotted = el.querySelectorAll('path[d$="5h2"]');
  animate(dotted, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin });

  const dashed = el.querySelectorAll('path[d$="12h6"]');
  animate(dashed, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: τmin });

  const solid = el.querySelector('path[d="M3 19h18"]');
  animate(solid, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: τmid });
}
