import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const shaved = el.querySelector('path[d^="M11 13h3"]');
  animate(shaved, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const glints = el.querySelectorAll('path[d="M19 13h2"], path[d="m18 17 1.5 1.5"], path[d="M14 19v2"]');
  animate(glints, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
}
