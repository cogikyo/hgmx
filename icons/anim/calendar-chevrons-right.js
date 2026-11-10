import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const page = el.querySelector('path[d^="M21 11.5"]');
  animate(page, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  const chevrons = [el.querySelector('path[d^="m19 21"]'), el.querySelector('path[d^="m13 21"]')];
  animate(chevrons, { translateX: [0, 2, -1, 0] }, { duration: 0.3, delay: stagger(0.12, { startDelay: τmin }) });
}
