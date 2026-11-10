import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const ring = el.querySelector("circle");
  animate(ring, { pathLength: [0, 1], opacity: [0.5, 1], rotate: [270, 90] }, { duration: τmid });

  const c = el.querySelector('path[d^="M15 9.4"]');
  animate(c, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
  const bar = el.querySelector('path[d="M7 12h5"]');
  animate(bar, { scaleX: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
