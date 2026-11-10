import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const rind = el.querySelectorAll('path[d^="M21.66"], path[d^="M19.65"]');
  animate(rind, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const segments = el.querySelectorAll('path[d^="m14 10"], path[d^="M14 17.85"]');
  animate(segments, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
