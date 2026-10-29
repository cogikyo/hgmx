import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const branches = el.querySelectorAll('path[d^="M3"]');
  animate(branches, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  const items = [el.querySelector('path[d="M8 5h13"]'), el.querySelector('path[d="M13 12h8"]'), el.querySelector('path[d="M13 19h8"]')];
  animate(items, { translateX: [-3, 0], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.06, { startDelay: τmin }) });
}
